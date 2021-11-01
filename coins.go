package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Coin struct {
	Name      string  `yaml:"name"`
	Amount    float64 `yaml:"amount"`
	TotalCost float64 `yaml:"totalCost"`
}

type CoinsToGetRidOf struct {
	Coin []string `yaml:"coins_to_get_want_to_get_rid_of"`
}

type Coins struct {
	Coins []Coin `yaml:"coins"`
}

func loadYamlForGetRidOf(fileName string) CoinsToGetRidOf {
	var config CoinsToGetRidOf
	source, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return config
}

func loadYamlForPortfolio(fileName string) Coins {
	var config Coins
	source, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return config

}
