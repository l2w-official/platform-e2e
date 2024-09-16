package internal

import (
	"platform_e2e/pkg/client/model"
	"platform_e2e/pkg/transport"
)

type APIClient struct {
	URL string
}

func (cli *APIClient) sendRequest(method, path string, req any, credentials *model.UserCredentials, res any, opts ...transport.RequestOpt) error {
	opts = append([]transport.RequestOpt{transport.WithBody(req), transport.WithContentType("application/ld+json")}, opts...)
	if credentials != nil {
		opts = append(opts, transport.WithCredentials(*credentials))
	}

	r, err := transport.NewRequest(method, cli.URL+path, opts...)
	if err != nil {
		return err
	}

	if _, err = r.Send(res); err != nil {
		return err
	}

	return nil
}
