package main

import (
	"log"
	"messagequeue/internal/config"
	"messagequeue/internal/generator"
	"messagequeue/internal/logging"
	"messagequeue/internal/sender"
)

const envPath = "config/.env"
const configPath = "config/config.yml"

func main() {
	// Парсинг из .env:
	// 		URL сохранения фактов,
	//		URL получения фактов,
	//		Bearer token
	//		auth_user_id
	// Парсинг из config.yml:
	//		настройки логгера,
	//		формат дат для проставления текущей даты,
	//		пресеты генератора фактов (часть полей и число требуемых фактов)
	config, err := config.LoadConfig(envPath, configPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Config loaded")

	// Формирование логгера на основе конфига
	// Переношу свой из проекта в проект :)
	// Можно заменить на общепроектный логгер
	logger, err := logging.NewLogrusLogger(&config.Logging)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("Logger configured")

	// Поля у отправлялки приватные, поэтому проставляются через конструктор
	sender := sender.NewFactSender(&config.App, &logger)

	// Генерация фактов
	facts := generator.GenerateFacts(&config.App)

	// Обработка
	sender.Process(facts)
}
