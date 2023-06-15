package util

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	clog "github.com/paceew/go-redisson/pkg/log"
)

func Recover(logger ...clog.FieldsLogger) {
	if err := recover(); err != nil {
		fmt.Printf("%s\n", err)
		fmt.Printf("%s\n", Stack())
		if len(logger) > 0 && logger[0] != nil {
			logger[0].Errorf("%s\n", err)
			logger[0].Errorf("%s\n", Stack())
		}
	}
}

func Stack() string {
	buf := make([]byte, 0, 2048)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

func GenUUID() string {
	return uuid.NewString()
}

func Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func LocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func PId() string {
	return strconv.Itoa(os.Getpid())
}

func GoId() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	return strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
}
