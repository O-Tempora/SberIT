package main

import (
	"flag"
	"fmt"
	"io"
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

	wr := getLoggerWriter()
	s := server.InitServer(cf).
		WithLogger(wr).
		WithDb(cf.DbHost, cf.DbBase, cf.DbPort)
	s.InitRouter()
	fmt.Println(s.Router)

	connectionInfo := fmt.Sprintf("%s:%d", cf.Host, cf.Port)
	s.Logger.Info().Msgf("Server starts on %s", connectionInfo)
	if err = http.ListenAndServe(connectionInfo, s); err != nil {
		s.Logger.Error().Msgf("Server start error: %s", connectionInfo)
	}
}

func getLoggerWriter() io.Writer {
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		return io.MultiWriter(os.Stdout)
	}
	file, err := os.OpenFile("logs/logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return io.MultiWriter(os.Stdout)
	}
	return io.MultiWriter(os.Stdout, file)
}
