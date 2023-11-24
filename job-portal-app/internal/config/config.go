package config

import (
	"log"

	env "github.com/Netflix/go-env"
)

var cfg Config

type Config struct {
	AppConfig   AppConfig
	DBConfig    DBConfig
	RedisConfig RedisConfig
	AuthConfig  AuthConfig
	MailConfig  MailConfig
}

type AppConfig struct {
	Host string `env:"APP_HOST"`
	Port string `env:"APP_PORT,required=true"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST"`
	Port     string `env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
}

type AuthConfig struct {
	PublicKey  string `env:"PUBLIC_KEY"`
	PrivateKey string `env:"PRIVATE_KEY"`
}

type MailConfig struct {
	Port int `env:"MAIL_PORT"`
}

type DBConfig struct {
	Host     string `env:"POSTGRES_HOST"`
	UserName string `env:"POSTGRES_USERNAME"`
	Port     string `env:"POSTGRES_PORT"`
	DbName   string `env:"POSTGRES_DB"`
	Password string `env:"POSTGRES_PASSWORD"`
	SslMode  string `env:"POSTGRES_SSL_MODE"`
	TimeZone string `env:"POSTGRES_TIME_ZONE"`
}

func init() {
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		log.Panic(err)
	}
}

func GetConfig() Config {
	return cfg
}
