package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func setupWebserver() {
	// Register prom metrics path in http serv
	httpMux := http.NewServeMux()
	httpMux.Handle("/metrics", promhttp.InstrumentMetricHandler(
		rConf.registry,
		promhttp.HandlerFor(rConf.registry, promhttp.HandlerOpts{}),
	))

	// Init & start serv
	rConf.httpServ = &http.Server{
		Addr:    fmt.Sprintf(":%d", rConf.httpServerPort),
		Handler: httpMux,
	}
	go func() {
		log.Infof("> Starting HTTP server at %s\n", rConf.httpServ.Addr)
		err := rConf.httpServ.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Errorf("HTTP Server errored out %v", err)
		}
	}()

}
