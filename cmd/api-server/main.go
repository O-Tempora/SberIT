package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/O-Tempora/SberIT/config"
	"github.com/O-Tempora/SberIT/internal/server"
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

	cf := config.Config{}
	err = yaml.Unmarshal(bytes, &cf)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err = os.MkdirAll("logs", os.ModePerm); err != nil {
		log.Fatal(err.Error())
	}
	file, err := os.OpenFile("logs/logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}

	s := server.InitServer(cf).
		WithLogger(os.Stdout, file).
		WithDb(cf.DbHost, cf.DbBase, cf.DbPort)
	s.Logger.Info().Msgf("Config: %+v", cf)
	s.Logger.Info().Msgf("Server: %+v", s)

	connectionInfo := fmt.Sprintf("%s:%d", cf.Host, cf.Port)
	s.Logger.Info().Msgf("Server starts on %s", connectionInfo)
	if err = http.ListenAndServe(connectionInfo, s); err != nil {
		s.Logger.Error().Msgf("Server start error: %s", connectionInfo)
	}
}
