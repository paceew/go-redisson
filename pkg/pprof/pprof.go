package pprof

import (
	"fmt"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// PprofLisent 如果pprofAddr不为空则开启pprof
//
// e.g : PprofLisent("127.0.0.1:16060")
func PprofLisent(pprofAddr string) {
	if pprofAddr != "" {
		go listenPPROF(pprofAddr)
	}
}

func listenPPROF(addr string) {
	routerPProf := gin.New()
	pprof.Register(routerPProf)
	fmt.Printf("pprof open on " + addr)
	if err := routerPProf.Run(addr); err != nil {
		fmt.Printf("pprof open err:%s", err.Error())
	}
}
