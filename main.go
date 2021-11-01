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
	vectors        map[string]*prometheus.GaugeVec
	httpServerPort uint
	httpServ       *http.Server
	updateInterval time.Duration
	debug          bool
	configFile     string
}

var rConf = runtimeConfStruct{
	httpServerPort: 9101,
	httpServ:       nil,
	vectors:        make(map[string]*prometheus.GaugeVec),
	registry:       prometheus.NewRegistry(),
	updateInterval: 50,
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

	coins := loadYamlForPortfolio(rConf.configFile)
	coinsToGetRidOf := loadYamlForGetRidOf(rConf.configFile)

	// Init Prometheus Gauge Vectors

	gaugeNames := []string{
		"amount", "total_costs", "get_rid_of",
	}
	for _, name := range gaugeNames {
		rConf.vectors[name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "crypto_balances",
			Name:      name,
			Help:      fmt.Sprintf(name),
		}, []string{"symbol"})
		rConf.registry.MustRegister(rConf.vectors[name])
	}

	// Regular loop operations below
	ticker := time.NewTicker(rConf.updateInterval)
	for {
		log.Debug("> Updating....\n")

		for _, item := range coinsToGetRidOf.Coin {
			rConf.vectors["get_rid_off"].WithLabelValues(item).Set(1)

		}
		for _, item := range coins.Coins {
			rConf.vectors["amount"].WithLabelValues(item.Name).Set(item.Amount)
			rConf.vectors["total_costs"].WithLabelValues(item.Name).Set(item.TotalCost)

		}
		<-ticker.C
	}
}
