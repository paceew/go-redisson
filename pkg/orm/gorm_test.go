package orm

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/paceew/go-redisson/pkg/config"
)

func TestGConfig(T *testing.T) {
	config.NewVipConfig("./test.yaml")
	dispatch_time := config.VipCfg().GetStringSlice("logger.dispatch_time")
	now := time.Now().Format("15:04")
	candis := false
	for _, v := range dispatch_time {
		ss := strings.Split(v, "-")
		if len(ss) == 2 {
			btime := ss[0]
			etime := ss[1]
			if btime < now {
				fmt.Println("yes")
			}

			if btime <= now && now <= etime {
				candis = true
				break
			}
		} else {
			fmt.Println("nono")
		}
	}
	fmt.Println(candis)
	// if dispatch_time == nil {
	// 	fmt.Println("nil")
	// }
	// if len(dispatch_time) == 0 {
	// 	fmt.Println("0")
	// }
	for k, v := range dispatch_time {
		fmt.Println(k, v, "22")
	}
	fmt.Println(dispatch_time)
	AutoInitGormWithConfig()
}
