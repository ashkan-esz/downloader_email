package configs

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type ConfigStruct struct {
	MainServerAddress  string
	RabbitMqUrl        string
	CorsAllowedOrigins []string
	MailServerHost     string
	MailServerPort     int
	MailServerUsername string
	MailServerPassword string
	UserSessionPage    string
}

var configs = ConfigStruct{}

func GetConfigs() ConfigStruct {
	return configs
}

func LoadEnvVariables() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	configs.MainServerAddress = os.Getenv("MAIN_SERVER_ADDRESS")
	configs.RabbitMqUrl = os.Getenv("RABBITMQ_URL")
	configs.CorsAllowedOrigins = strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), "---")
	for i := range configs.CorsAllowedOrigins {
		configs.CorsAllowedOrigins[i] = strings.TrimSpace(configs.CorsAllowedOrigins[i])
	}
	configs.MailServerHost = os.Getenv("MAILSERVER_HOST")
	configs.MailServerPort, _ = strconv.Atoi(os.Getenv("MAILSERVER_PORT"))
	configs.MailServerUsername = os.Getenv("MAILSERVER_USERNAME")
	configs.MailServerPassword = os.Getenv("MAILSERVER_PASSWORD")
	configs.UserSessionPage = os.Getenv("USER_SESSION_PAGE")
}
