package models

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SecretKey   string `yaml:"SecretKey"`
	DataBaseURI string `yaml:"DataBaseURI"`
}

func (c *Config) InitConfig() *Config {
	file, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(file, c)
	if err != nil {
		log.Fatal(err)
	}

	if len(c.SecretKey) == 0 {
		log.Fatal("SecretKey must be set in config.yaml")
	}

	if c.DataBaseURI == "" {
		log.Fatal("DataBaseURI must be set in config.yaml")
	}

	return c
}
