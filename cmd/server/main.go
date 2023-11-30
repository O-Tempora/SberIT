package main

import (
	"flag"
	"log"
	"os"

	"github.com/O-Tempora/SberIT/config"
	"gopkg.in/yaml.v3"
)

const defaultConfig = "config/default.yaml"

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", defaultConfig, "Path to config file")
}

func main() {
	flag.Parse()
	if configPath == "" {
		configPath = defaultConfig
	}

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	cf := &config.Config{}
	err = yaml.Unmarshal(bytes, cf)
	if err != nil {
		log.Fatal(err.Error())
	}

}
