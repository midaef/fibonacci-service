package main

import (
	"flag"
	"github.com/midaef/fibonacci-service/app"
	"github.com/midaef/fibonacci-service/config"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

const prodConfigPath = "./config/prod-config.yaml"

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "./config/prod-config.yaml", "path to config")
}

func main() {
	flag.Parse()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Fatal("failed config file reading`",
			zap.String("error", err.Error()),
		)
	}

	var config config.Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		logger.Fatal("failed config unmarshalling",
			zap.String("error", err.Error()),
		)
	}

	app := app.NewApp(logger, &config)
	err = app.StartApp()
	if err != nil {
		logger.Fatal("failed start app",
			zap.String("error", err.Error()),
		)
	}
}
