package main

import (
	"log"
	"messagequeue/internal/config"
	"messagequeue/internal/generator"
	"messagequeue/internal/logging"
	"messagequeue/internal/pkg/entities"
	"messagequeue/internal/sender"
	"time"
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

	preset := entities.GeneratorPreset{
		PeriodStart: time.Date(2005, 8, 1, 0, 0, 0, 0, nil).String(),
		PeriodEnd:   time.Date(2005, 8, 30, 0, 0, 0, 0, nil).String(),
		PeriodKey:   "month",
		MoID:        "999999",
		IsPlan:      "0",
		AuthUserID:  config.App.AuthUserID,
	}

	facts := generator.GenerateFacts(2000, preset)

	sender.Process(facts)
}
