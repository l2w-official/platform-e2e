package main_suite_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type request struct {
	req *http.Request
}

type requestOpt func(r *request) error

func newRequest(method string, url string, opts ...requestOpt) (req request, err error) {
	if req.req, err = http.NewRequest(method, url, nil); err != nil {
		return req, err
	}

	for _, opt := range opts {
		if err = opt(&req); err != nil {
			return req, err
		}
	}

	return req, nil
}

func withBody(v any) requestOpt {
	return func(r *request) error {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}

		fmt.Println("### body", string(b))
		r.req.Body = io.NopCloser(bytes.NewReader(b))
		return nil
	}
}

func withHeader(key, value string) requestOpt {
	return func(r *request) error {
		r.req.Header.Add(key, value)
		return nil
	}
}

func withCredentials(credentials userCredentials) requestOpt {
	return func(r *request) error {
		r.req.Header.Set("Authorization", "Bearer "+credentials.accessToken)
		r.req.Header.Set("x-context-token", credentials.contextToken)
		return nil
	}
}

func withContentType(contentType string) requestOpt {
	return func(r *request) error {
		r.req.Header.Set("Content-Type", contentType)
		return nil
	}
}

func withQueryParam(key, value string) requestOpt {
	return func(r *request) error {
		q := r.req.URL.Query()
		q.Add(key, value)
		r.req.URL.RawQuery = q.Encode()
		return nil
	}
}

type httpError struct {
	message string
	code    int
}

func (e *httpError) Error() string {
	return fmt.Sprintf("Error %d: %s.",
		e.code,
		e.message,
	)
}

func newHttpError(message string, code int) *httpError {
	return &httpError{message, code}
}

func (r *request) send(v any) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(r.req)
	if err != nil {
		fmt.Println("### err", err.Error())

		return resp, err
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}

	fmt.Println("### resp", string(bytes))

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return nil, newHttpError(string(bytes), resp.StatusCode)
	}

	if v == nil {
		return resp, nil
	}

	if err = json.Unmarshal(bytes, v); err != nil {
		return resp, err
	}

	return resp, nil
}

func executeHttpRequest(req *http.Request, v any) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("### err", err.Error())

		return resp, err
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}

	fmt.Println("### resp", string(bytes))

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("Error %d: %s", resp.StatusCode, string(bytes))
	}

	if v == nil {
		return resp, nil
	}

	if err = json.Unmarshal(bytes, v); err != nil {
		return resp, err
	}

	return resp, nil
}

func login(orgID, username, password string) (string, error) {
	form := url.Values{}
	form.Add("client_id", "l2w-app")
	form.Add("username", username)
	form.Add("password", password)
	form.Add("grant_type", "password")

	url := fmt.Sprintf("http://localhost:8080/realms/%s/protocol/openid-connect/token", orgID)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var loginResponse struct {
		AccessToken string `json:"access_token"`
	}

	if _, err = executeHttpRequest(req, &loginResponse); err != nil {
		return "", err
	}

	return loginResponse.AccessToken, nil
}

func userContext(orgID, token string, acceptTerms bool) (string, error) {
	reqBody := map[string]any{
		"org_id": orgID,
	}
	if acceptTerms {
		for _, t := range []string{"accept_privacy_policy", "accept_terms_and_conditions"} {
			reqBody[t] = true
		}
	}

	url := "http://localhost:8050/v1/contexts"
	r, err := newRequest(http.MethodPost, url, withBody(reqBody), withHeader("Authorization", fmt.Sprintf("Bearer %s", token)))
	if err != nil {
		return "", err
	}

	var contextResponse struct {
		Token string `json:"token"`
	}

	if _, err = r.send(&contextResponse); err != nil {
		return "", err
	}

	return contextResponse.Token, nil
}

type organization struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	IDPType     string `json:"idpType"`
	IDPGroupID  string `json:"idpGroupId"`
	IDPClientID string `json:"idpClientId"`
	Status      string `json:"statis"`
}

type orgAttribute struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Type             string             `json:"types"`
	Status           string             `json:"status"`
	AttributeOptions []*attributeOption `json:"attributeOptions"`
}

func orgAttributes(idpToken, contextToken string) ([]*orgAttribute, error) {
	r, err := newRequest(http.MethodGet, "http://localhost:8050/v1/attributes?page=1&itemsPerPage=10&matchStatus=ACTIVE&matchEditable=true",
		withHeader("Authorization", fmt.Sprintf("Bearer %s", idpToken)),
		withHeader("x-context-token", contextToken))

	if err != nil {
		return nil, err
	}

	var orgAttributesResp struct {
		Attributes []*orgAttribute `json:"attributes"`
	}
	if _, err = r.send(&orgAttributesResp); err != nil {
		return nil, err
	}

	return orgAttributesResp.Attributes, nil
}
