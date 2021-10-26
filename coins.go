package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Name   string  `yaml:"name"`
	Amount float64 `yaml:"amount"`
}

type Coins struct {
	Coins []Config `yaml:"coins"`
}

func loadYaml(fileName string, config Coins) Coins {

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
