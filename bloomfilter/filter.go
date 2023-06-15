package bloomfilter

import (
	"errors"
	"fmt"
	"github.com/paceew/go-redisson/pkg/log"
	"math"
	"strconv"
	"time"
)

var (
	ErrGreaterFalseProbability  = errors.New("filter false probability can't be greater than 1")
	ErrNegativeFalseProbability = errors.New("filter false probability can't be negative")
	ErrSizeZero                 = errors.New("filter calculated size is zero")
	ErrSizeTooBig               = fmt.Errorf("filter calculated size is greater than %d", RedisMaxLength)
	ErrNotInit                  = errors.New("filter is not init")
	ErrConfigUnequal            = errors.New("filter config unequal")
)

type GFilter struct {
	gBloomFilter       *GBloomFilter
	key                string
	configKey          string
	size               uint64  // bitmap size
	hashIterations     uint64  // 哈希次数
	expectedInsertions uint64  // 预期插入量
	falseProbability   float64 // 预期错误率
	hashfunc           HashFunc
	checkConfig        bool // 操作前，是否需要对配置进行检查
	logger             log.FieldsLogger
}

type HashFunc func(data []byte, hashIterations uint64, size uint64) []uint64

// TryInit 尝试初始化，如果已经初始化(对于相同的key包括已经在其他机器初始化，设置了redis key config)则读取redis配置信息覆盖当前的配置
//
// @expectedInsertions 预期插入量
// @falseProbability 预期错误率
func (gf *GFilter) TryInit(expectedInsertions uint64, falseProbability float64) error {
	if ok, err := gf.readConfig(); err != nil {
		return err
	} else if !ok {
		if falseProbability > 1 {
			return ErrGreaterFalseProbability
		} else if falseProbability < 0 {
			return ErrNegativeFalseProbability
		}

		size := gf.optimalNumOfBits(expectedInsertions, falseProbability)
		if size == 0 {
			return ErrSizeZero
		} else if size > RedisMaxLength {
			return ErrSizeTooBig
		}

		hashIterations := gf.optimalNumOfHashFunctions(expectedInsertions, size)

		gf.size = size
		gf.hashIterations = hashIterations
		gf.expectedInsertions = expectedInsertions
		gf.falseProbability = falseProbability
		err := gf.writeConfig()
		if err != nil {
			return err
		}
	}

	return nil
}

func (gf *GFilter) optimalNumOfBits(expectedInsertions uint64, falseProbability float64) (size uint64) {
	return uint64(-1 * float64(expectedInsertions) * math.Log(falseProbability) / math.Pow(math.Log(2), 2))
}

func (gf *GFilter) optimalNumOfHashFunctions(expectedInsertions, size uint64) (hashIterations uint64) {
	return uint64(math.Max(1, math.Ceil(float64(size)/float64(expectedInsertions)*math.Log(2))))
}

func (gf *GFilter) readConfig() (bool, error) {
	config, err := gf.gBloomFilter.redisOper.HMGet(gf.configKey, "size", "hashIterations", "expectedInsertions", "falseProbability")
	if err != nil {
		return false, err
	}
	if len(config) != 0 {
		if size, ok := config["size"]; ok {
			if ss, err := strconv.ParseUint(size, 10, 64); err == nil {
				gf.size = ss
			}
		}

		if hashIterations, ok := config["hashIterations"]; ok {
			if hi, err := strconv.ParseUint(hashIterations, 10, 64); err == nil {
				gf.hashIterations = hi
			}
		}

		if expectedInsertions, ok := config["expectedInsertions"]; ok {
			if ei, err := strconv.ParseUint(expectedInsertions, 10, 64); err == nil {
				gf.expectedInsertions = ei
			}
		}

		if falseProbability, ok := config["falseProbability"]; ok {
			if fp, err := strconv.ParseFloat(falseProbability, 64); err == nil {
				gf.falseProbability = fp
			}
		}

		return true, nil
	}

	return false, nil
}

func (gf *GFilter) writeConfig() error {
	sizestr := strconv.FormatUint(gf.size, 10)
	hashIterationsstr := strconv.FormatUint(gf.hashIterations, 10)
	expectedInsertionsstr := strconv.FormatUint(gf.expectedInsertions, 10)
	falseProbabilitystr := strconv.FormatFloat(gf.falseProbability, 'f', 10, 64)
	err := gf.gBloomFilter.redisOper.HMSet(gf.configKey, "size", sizestr, "hashIterations", hashIterationsstr,
		"expectedInsertions", expectedInsertionsstr, "falseProbability", falseProbabilitystr)
	return err
}

func (gf *GFilter) configCheck() error {
	if gf.size == 0 || gf.hashIterations == 0 {
		if ok, err := gf.readConfig(); err != nil {
			return err
		} else if !ok {
			return ErrNotInit
		}
	} else if gf.checkConfig {
		if config, err := gf.gBloomFilter.redisOper.HMGet(gf.configKey, "size", "hashIterations"); err != nil {
			return err
		} else {
			var size uint64
			var hashIterations uint64
			if sizestr, ok := config["size"]; ok {
				size, _ = strconv.ParseUint(sizestr, 10, 64)
			}
			if hashIterationsstr, ok := config["hashIterations"]; ok {
				hashIterations, _ = strconv.ParseUint(hashIterationsstr, 10, 64)
			}

			if gf.size != size || gf.hashIterations != hashIterations {
				return ErrConfigUnequal
			}
		}
	}
	return nil
}

func (gf *GFilter) Add(obj interface{}) error {
	if err := gf.configCheck(); err != nil {
		return err
	}

	objstr := fmt.Sprintf("%s", obj)
	indexs := gf.hashfunc([]byte(objstr), gf.hashIterations, gf.size)
	for _, v := range indexs {
		if err := gf.gBloomFilter.redisOper.SetBit(gf.key, v%gf.size, 1); err != nil {
			return err
		}
	}
	return nil
}

func (gf *GFilter) Contains(obj interface{}) (bool, error) {
	if err := gf.configCheck(); err != nil {
		return false, err
	}

	objstr := fmt.Sprintf("%s", obj)
	indexs := gf.hashfunc([]byte(objstr), gf.hashIterations, gf.size)
	for _, v := range indexs {
		if flag, err := gf.gBloomFilter.redisOper.GetBit(gf.key, v%gf.size); err != nil {
			return false, err
		} else if flag == 0 {
			return false, nil
		}
	}
	return true, nil
}

func (gf *GFilter) Expired(second time.Duration) (int64, error) {
	if err := gf.configCheck(); err != nil {
		return 0, err
	}

	gf.gBloomFilter.redisOper.Expire(gf.configKey, second)
	return gf.gBloomFilter.redisOper.Expire(gf.key, second)
}

func (gf *GFilter) Del() error {
	if err := gf.configCheck(); err != nil {
		return err
	}

	gf.gBloomFilter.redisOper.Del(gf.configKey)
	return gf.gBloomFilter.redisOper.Del(gf.key)
}
