package version

import (
	"fmt"
	"os"
)

var (
	_COMMIT_  = "unknown"
	_BTIME_   = "unknown"
	_VERSION_ = "unknown"
	_BRANCH_  = "unknown"
	PROGRAME_ = "unknown"
)

// VersionController 版本号控制，输出Version信息到标准输出并且正常结束程序
// isVersion,是否输出版本号，常见的一种用法是:
//
//isVersion := flag.Bool("v", false, "print version and exit"),
//VersionController(programe,isVersion)
func VersionController(isVersion *bool) {
	if *isVersion {
		fmt.Printf("%s\n", FetchVersionStr())
		os.Exit(0)
	}
}

func FetchVersionStr() string {
	return fmt.Sprintf("programe:%s,version:%s,branch:%s,commit:%s,build time:%s", PROGRAME_, _VERSION_, _BRANCH_, _COMMIT_, _BTIME_)
}
