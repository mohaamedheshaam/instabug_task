package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	MySQL         MySQLConfig
	Redis         RedisConfig
	RabbitMQ      RabbitMQConfig
	Elasticsearch ElasticsearchConfig
}

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type RedisConfig struct {
	Host string
	Port string
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type ElasticsearchConfig struct {
	URL string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		MySQL: MySQLConfig{
			Host:     viper.GetString("MYSQL_HOST"),
			Port:     viper.GetString("MYSQL_PORT"),
			User:     viper.GetString("MYSQL_USER"),
			Password: viper.GetString("MYSQL_PASSWORD"),
			Database: viper.GetString("MYSQL_DATABASE"),
		},
		Redis: RedisConfig{
			Host: viper.GetString("REDIS_HOST"),
			Port: viper.GetString("REDIS_PORT"),
		},
		RabbitMQ: RabbitMQConfig{
			Host:     viper.GetString("RABBITMQ_HOST"),
			Port:     viper.GetString("RABBITMQ_PORT"),
			User:     viper.GetString("RABBITMQ_USER"),
			Password: viper.GetString("RABBITMQ_PASSWORD"),
		},
		Elasticsearch: ElasticsearchConfig{
			URL: viper.GetString("ELASTICSEARCH_URL"),
		},
	}, nil
}
