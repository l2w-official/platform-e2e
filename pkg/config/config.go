package config

import "os"

type Config struct {
	SuperAdminOrgID    string
	SuperAdminEmail    string
	SuperAdminPassword string

	ApiGatewayURL string

	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
}

var config *Config

func Get() *Config {
	if config == nil {
		loadConfig()
	}

	return config
}

func loadConfig() {
	if config != nil {
		return
	}

	config = &Config{
		SuperAdminOrgID:    os.Getenv("SUPER_ADMIN_ORG_ID"),
		SuperAdminEmail:    os.Getenv("SUPER_ADMIN_EMAIL"),
		SuperAdminPassword: os.Getenv("SUPER_ADMIN_PASSWORD"),
		ApiGatewayURL:      os.Getenv("API_GATEWAY_URL"),
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             os.Getenv("DB_PORT"),
	}
}
