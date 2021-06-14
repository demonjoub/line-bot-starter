package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/demonjoub/chatbot/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Configuration struct {
	Line        LineApi   `mapstructure:"lineApi"`
	App         AppConfig `mapstructure:"app"`
	Environment string    `mapstructure:"environment"`
}

type LineApi struct {
	ChannelID          int    `mapstructure:"channelId"`
	ChannelSecret      string `mapstructure:"channelSecret"`
	UserId             string `mapstructure:"userId"`
	ChannelAccessToken string `mapstructure:"channelAccessToken"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type Config struct {
	Environment string
	App         AppConfig
	Line        LineApi
}

func NewConfig() *Config {
	configPath, ok := os.LookupEnv("API_CONFIG_PATH")
	if !ok {
		logger.Info("API_CONFIG_PATH is not found, use default config file")
		configPath = "."
	}
	configName, ok := os.LookupEnv("API_CONFIG_NAME")
	if !ok {
		logger.Info("API_CONFIG_NAME is not found, use default config name")
		configName = "config"
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)

	if err := viper.ReadInConfig(); err != nil {
		logger.Error(err.Error(), zap.String("config", "config file not found"))
		return &Config{}
	}
	viper.AutomaticEnv()
	viper.WatchConfig()

	var configuration Configuration
	err := viper.Unmarshal(&configuration)
	if err != nil {
		logger.Error(err.Error(), zap.String("config", "viper Unmarshal error"))
		return &Config{}
	}

	logger.Info("LINE API", zap.String("ChannelID", fmt.Sprintf("%d", configuration.Line.ChannelID)))
	logger.Info("LINE API", zap.String("UserId", configuration.Line.UserId))

	logger.Info("environment", zap.String("env", configuration.Environment))

	return &Config{
		App:         configuration.App,
		Line:        configuration.Line,
		Environment: configuration.Environment,
	}
}

func GetSecretValue() {
	for _, value := range os.Environ() {
		pair := strings.SplitN(value, "=", 2)
		if strings.Contains(pair[0], "SECRET_") == true {
			keys := strings.Replace(pair[0], "SECRET_", "secrets.", -1)
			keys = strings.Replace(keys, "_", ".", -1)
			newKey := strings.Trim(keys, " ")
			newValue := strings.Trim(pair[1], " ")
			viper.Set(newKey, newValue)
		}
	}
}
