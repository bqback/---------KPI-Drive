package config

import (
	"os"

	"messagequeue/internal/apperrors"
	"messagequeue/internal/pkg/entities"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

// Config общее хранилище конфига
type Config struct {
	App     AppConfig     `yaml:"app"`
	Logging LoggingConfig `yaml:"logging"`
}

// AppConfig конфиг приложения, собираемый из .env и config.yml
type AppConfig struct {
	SaveURL         string                   `yaml:"-"`
	Token           string                   `yaml:"-"`
	AuthUserID      string                   `yaml:"-"`
	MaxRequests     int                      `yaml:"max_requests"`
	DateFormat      string                   `yaml:"date_layout"`
	GeneratorPreset entities.GeneratorPreset `yaml:"generator_preset"`
}

// LoggingConfig конфиг логгера, собираемый из .env и config.yml
type LoggingConfig struct {
	Level                  string `yaml:"level"`
	DisableTimestamp       bool   `yaml:"disable_timestamp"`
	FullTimestamp          bool   `yaml:"full_timestamp"`
	DisableLevelTruncation bool   `yaml:"disable_level_truncation"`
	LevelBasedReport       bool   `yaml:"level_based_report"`
	ReportCaller           bool   `yaml:"report_caller"`
}

// LoadConfig парсит .env и config.yml по указанным путям в конфиг
func LoadConfig(envPath string, configPath string) (*Config, error) {
	var config Config

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	err = godotenv.Load(envPath)

	if err != nil {
		return nil, apperrors.ErrEnvNotFound
	}

	// У каждой переменной из .env свой геттер, чтобы отдавать в нём кастомную ошибку
	// В принципе можно и заменить содержимым функции
	saveURL, err := GetSaveURL()
	if err != nil {
		return nil, err
	}

	token, err := GetBearerToken()
	if err != nil {
		return nil, err
	}

	id, err := GetAuthUserID()
	if err != nil {
		return nil, err
	}

	config.App.SaveURL = saveURL
	config.App.Token = token
	config.App.AuthUserID = id

	return &config, nil
}

func GetSaveURL() (string, error) {
	url, ok := os.LookupEnv("SAVE_URL")
	if !ok {
		return url, apperrors.ErrSaveURLMissing
	}
	return url, nil
}

func GetBearerToken() (string, error) {
	token, ok := os.LookupEnv("TOKEN")
	if !ok {
		return token, apperrors.ErrBearerTokenMissing
	}
	return token, nil
}

func GetAuthUserID() (string, error) {
	id, ok := os.LookupEnv("USER_ID")
	if !ok {
		return id, apperrors.ErrAuthUserIDMissing
	}
	return id, nil
}
