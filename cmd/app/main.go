package main

import (
	"log"
	"messagequeue/internal/config"
	"messagequeue/internal/logging"
	"messagequeue/internal/sender"
)

const envPath = "config/.env"
const configPath = "config/config.yml"

func main() {
	config, err := config.LoadConfig(envPath, configPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Config loaded")

	logger, err := logging.NewLogrusLogger(config.Logging)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("Logger configured")

	sender := sender.NewFactSender(config.App, &logger)
}
