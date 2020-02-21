package main

import (
	"btc-project/conf"
	"btc-project/internal/server/http"
	"btc-project/internal/service"
	"btc-project/library/log"
	"context"
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	// init timezone as +0
	if err := os.Setenv("TZ", "Greenwich"); err != nil {
		panic(err)
	}

	// init conf
	if err := conf.Init(); err != nil {
		panic(err)
	}

	// init log
	if conf.Conf.Constants.IsDev {
		log.Init(nil)
	} else {

		log.Init(&log.Config{
			Dir: "../logs",
		})

	}
	defer log.Close()

	log.Info("btc-project start")
	svc := service.New(conf.Conf)
	httpSrv := http.New(conf.Conf, svc)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
			if err := httpSrv.Shutdown(ctx); err != nil {
				log.Error("httpSrv.Shutdown error(%v)", err)
			}
			log.Info("btc-project exit")
			svc.Close()
			cancel()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
