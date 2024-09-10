package main_suite_test

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type keycloakCli struct {
	adminToken string
}

func newKeycloakCli() *keycloakCli {
	cli := &keycloakCli{}
	cli.refreshToken()
	return cli
}

func (cli *keycloakCli) refreshToken() {
	form := url.Values{}
	form.Add("client_id", "admin-cli")
	form.Add("username", "admin")
	form.Add("password", "admin")
	form.Add("grant_type", "password")

	url := fmt.Sprintf("http://localhost:8080/realms/%s/protocol/openid-connect/token", "master")
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal("could not get keycloak admin token")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var loginResponse struct {
		AccessToken string `json:"access_token"`
	}

	if _, err = executeHttpRequest(req, &loginResponse); err != nil {
		log.Fatal("could not get keycloak admin token")
	}

	cli.adminToken = loginResponse.AccessToken
}

func (cli *keycloakCli) deleteRealm(realmID string) error {
	r, err := newRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8080/admin/realms/%s", realmID),
		withHeader("Authorization", fmt.Sprintf("Bearer %s", cli.adminToken)))
	if err != nil {
		return err
	}

	_, err = r.send(nil)
	return err
}
