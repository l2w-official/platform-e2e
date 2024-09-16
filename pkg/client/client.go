package client

import (
	"platform_e2e/pkg/client/internal"
	"platform_e2e/pkg/config"
)

func New(cnf *config.Config) APIClient {
	return &internal.APIClient{
		URL: cnf.ApiGatewayURL,
	}
}
