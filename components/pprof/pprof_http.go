package pprof

import (
	"github.com/DeanThompson/ginpprof"

	"gitlab.yidian-inc.com/image/lego/components/httpserver"
)

func UseHttpPprof(server *httpserver.HttpServer) {
	//http server 设置
	ginpprof.Wrap(server.Engine)
}
