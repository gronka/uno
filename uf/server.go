package uf

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/fasthttp/router"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func cleanup() {
	csb.Close()
	csm.Close()
}

func InitAndListenAndServe(conf Config, router *router.Router) {
	if conf.ConnectCassandra == true {
		connectCassandra(true)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	reinit := flag.Bool("reinit", false, "does nothing rn")
	if *reinit == true {
		Debug("reinit flag currently does nothing")
	} else {
		//TODO: check that all errors are wrapped before coming here
		if conf.UfName == "public" && conf.Environment == "local" {
			Fatal(errors.Cause(fasthttp.ListenAndServeTLS(
				conf.BindTo,
				"./localhost+8.pem",
				"./localhost+8-key.pem",
				middlewareCors(router.Handler),
			)))
		} else {
			Fatal(errors.Cause(fasthttp.ListenAndServe(
				conf.BindTo,
				middlewareCors(router.Handler),
			)))
		}
	}

}
