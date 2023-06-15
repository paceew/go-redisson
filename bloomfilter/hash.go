package bloomfilter

import (
	"github.com/spaolacci/murmur3"
)

func HashMur3(data []byte, hashIterations uint64, size uint64) []uint64 {
	hasher := murmur3.New128()
	hasher.Write(data)
	v1, v2 := hasher.Sum128()

	indexes := make([]uint64, hashIterations)
	hash := v1
	for i := 0; i < int(hashIterations); i++ {
		indexes[i] = hash % size
		if i%2 == 0 {
			hash += v2
		} else {
			hash += v1
		}
	}

	return indexes
}
