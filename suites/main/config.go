package main_suite_test

import "os"

type cnf struct {
	superAdminOrgID    string
	superAdminEmail    string
	superAdminPassword string

	apiGatewayURL string

	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     string
}

var config *cnf

func loadConfig() {
	if config != nil {
		return
	}

	config = &cnf{
		superAdminOrgID:    os.Getenv("SUPER_ADMIN_ORG_ID"),
		superAdminEmail:    os.Getenv("SUPER_ADMIN_EMAIL"),
		superAdminPassword: os.Getenv("SUPER_ADMIN_PASSWORD"),
		apiGatewayURL:      os.Getenv("API_GATEWAY_URL"),
		dbUser:             os.Getenv("DB_USER"),
		dbPassword:         os.Getenv("DB_PASSWORD"),
		dbHost:             os.Getenv("DB_HOST"),
		dbPort:             os.Getenv("DB_PORT"),
	}
}
