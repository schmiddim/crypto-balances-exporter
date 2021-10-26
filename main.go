package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type runtimeConfStruct struct {
	registry       *prometheus.Registry
	balances       *prometheus.GaugeVec
	httpServerPort uint
	httpServ       *http.Server
	updateIval     time.Duration
	debug          bool
	configFile     string
}

var rConf = runtimeConfStruct{

	httpServerPort: 9101,
	httpServ:       nil,
	registry:       prometheus.NewRegistry(),
	updateIval:     50,
	balances:       nil,
	configFile:     "",
}

func initParams() {

	flag.UintVar(&rConf.httpServerPort, "httpServerPort", rConf.httpServerPort, "HTTP server port.")
	flag.BoolVar(&rConf.debug, "debug", false, "Set debug log level.")
	flag.StringVar(&rConf.configFile, "config", "config.yaml", "Path to configfile")
	flag.Parse()

	logLvl := log.InfoLevel
	if rConf.debug {
		logLvl = log.DebugLevel
	}
	log.SetLevel(logLvl)

}

func main() {
	initParams()
	var coins Coins

	coins = loadYaml(rConf.configFile, coins)

	fmt.Printf("coins:\n%v\n\n", coins)

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

	// Init Prometheus Gauge Vectors
	rConf.balances = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "fuckyou",
		Name:      "balance",
		Help:      fmt.Sprintf("Balance in account for assets"),
	},
		[]string{"symbol"},
	)
	rConf.registry.MustRegister(rConf.balances)

	// Regular loop operations below
	ticker := time.NewTicker(rConf.updateIval)
	for {
		log.Debug("> Updating....\n")

		for _, item := range coins.Coins {
			rConf.balances.WithLabelValues(item.Name).Set(item.Amount)

		}
		<-ticker.C
	}
}
