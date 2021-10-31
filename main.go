package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type runtimeConfStruct struct {
	registry       *prometheus.Registry
	amounts        *prometheus.GaugeVec
	totalCosts     *prometheus.GaugeVec
	httpServerPort uint
	httpServ       *http.Server
	updateInterval time.Duration
	debug          bool
	configFile     string
}

var rConf = runtimeConfStruct{

	httpServerPort: 9101,
	httpServ:       nil,
	registry:       prometheus.NewRegistry(),
	updateInterval: 50,
	amounts:        nil,
	totalCosts:     nil,
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
	setupWebserver()
	var coins Coins

	coins = loadYaml(rConf.configFile, coins)

	// Init Prometheus Gauge Vectors

	rConf.amounts = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "crypto_balances",
		Name:      "amount",
		Help:      fmt.Sprintf("Balance in account for assets"),
	}, []string{"symbol"})
	rConf.totalCosts = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "crypto_balances",
		Name:      "total_costs",
		Help:      fmt.Sprintf("total costs for the coins"),
	}, []string{"symbol"})
	rConf.registry.MustRegister(rConf.amounts)
	rConf.registry.MustRegister(rConf.totalCosts)

	// Regular loop operations below
	ticker := time.NewTicker(rConf.updateInterval)
	for {
		log.Debug("> Updating....\n")

		for _, item := range coins.Coins {
			rConf.amounts.WithLabelValues(item.Name).Set(item.Amount)
			rConf.totalCosts.WithLabelValues(item.Name).Set(item.TotalCost)
		}
		<-ticker.C
	}
}
