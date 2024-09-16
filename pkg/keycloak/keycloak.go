package keycloak

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"platform_e2e/pkg/transport"
)

type Client struct {
	adminToken string
}

func New() *Client {
	cli := &Client{}
	cli.refreshToken()
	return cli
}

func (cli *Client) refreshToken() {
	form := url.Values{}
	form.Add("client_id", "admin-cli")
	form.Add("username", "admin")
	form.Add("password", "admin")
	form.Add("grant_type", "password")

	url := fmt.Sprintf("http://localhost:8080/realms/%s/protocol/openid-connect/token", "master")

	var loginResponse struct {
		AccessToken string `json:"access_token"`
	}

	req, err := transport.NewRequest(http.MethodPost, url, transport.WithForm(form), transport.WithContentType("application/x-www-form-urlencoded"))
	if err != nil {
		log.Fatal("could not create keycloak token request")
	}

	if _, err = req.Send(&loginResponse); err != nil {
		log.Fatal("could not get keycloak admin token")
	}

	cli.adminToken = loginResponse.AccessToken
}

func (cli *Client) DeleteRealm(realmID string) error {
	r, err := transport.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8080/admin/realms/%s", realmID),
		transport.WithHeader("Authorization", fmt.Sprintf("Bearer %s", cli.adminToken)))
	if err != nil {
		return err
	}

	_, err = r.Send(nil)
	return err
}
