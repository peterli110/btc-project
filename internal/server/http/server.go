package http

import (
	"btc-project/conf"
	"btc-project/internal/service"
	"btc-project/library/ecode"
	"btc-project/library/log"
	bm "btc-project/library/net/http/blademaster"
	"net/http"
)

var (
	svc 				*service.Service
)


// New init server.
func New(c *conf.Config,s *service.Service) (engine *bm.Engine) {
	// init http engine
	svc = s
	engine = bm.DefaultServer(c.Server)

	initRouter(engine)

	if err := engine.Start(); err != nil {
		panic(err)
	}
	return
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/coin")
	{
		g.POST("/add", addCoin)
		g.POST("/search", queryCoin)
	}

}

func ping(ctx *bm.Context) {
	if err := svc.Ping(ctx); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}


func pong(ctx *bm.Context) {
	ctx.JSON(nil, ecode.InvalidParams)
}