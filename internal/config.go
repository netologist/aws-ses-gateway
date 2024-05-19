package internal

import (
	"log"
	"github.com/caarlos0/env/v7"
)

type ConfigType struct {
	Port     int    `env:"PORT" envDefault:"8081"`
	SmtpHost string `env:"SMTP_HOST" envDefault:"localhost"`
	SmtpPort int `env:"SMTP_PORT" envDefault:"25"`
	SmtpUser string `env:"SMTP_USER"`
	SmtpPass string `env:"SMTP_PASS"`
}

var Config ConfigType

func ReadConfigFromEnv() {
	Config = ConfigType{}
	if err := env.Parse(&Config); err != nil {
		log.Printf("%+v", err)
	}

	if Config.SmtpHost == "" {
		panic("SMTP_HOST not defined")
	}

	log.Printf("%+v", Config)
}
