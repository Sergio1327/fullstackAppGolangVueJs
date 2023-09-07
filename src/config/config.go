package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/configor"
)

// Config конфиг
type Config struct {
	Redis struct {
		CacheTime time.Duration `yaml:"cacheTime"`
	} `yaml:"redis"`
}

// NewConfig init and return project config
func NewConfig(confPath string) (Config, error) {
	var c = Config{}
	err := configor.Load(&c, confPath)
	return c, err
}

// PostgresURL connect to callcenter's postgres db
func (c *Config) PostgresURL() string {
	pgURL := os.Getenv("PG_URL")
	if pgURL == "" {
		log.Fatalln("отсутствует PG_URL")
	}

	return fmt.Sprintf("postgres://%s?sslmode=disable", pgURL)
}

// OracleConnectString connect to oracle server
func (c *Config) OracleConnectString() string {
	path := os.Getenv("ORACLE_CONNECT_FILE")
	if path == "" {
		log.Fatalln("отсутствует ORACLE_CONNECT_FILE")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("не удалось считать ORACLE_CONNECT_FILE")
	}
	return string(data)
}

// RedisConnect connect to redis server
func (c *Config) RedisConnect() string {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatalln("отсутствует REDIS_URL")
	}

	return fmt.Sprintf("redis://%s/", redisURL)
}

// RedisCacheLifetime cache life time
func (c *Config) RedisCacheLifetime() time.Duration {
	return c.Redis.CacheTime * time.Hour
}

// RabbitMQConnectURL подключение к rabbitmq
func (c *Config) RabbitMQConnectURL() string {
	rabbitURL := os.Getenv("RABBIT_URL")
	if rabbitURL == "" {
		log.Fatalln("отсутствует RABBIT_URL")
	}

	return fmt.Sprintf("amqp://%s/", rabbitURL)
}
