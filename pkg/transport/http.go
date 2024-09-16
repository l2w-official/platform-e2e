package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"platform_e2e/pkg/client/model"
	"strings"
)

type Request struct {
	req *http.Request
}

type RequestOpt func(r *Request) error

func NewRequest(method string, url string, opts ...RequestOpt) (req Request, err error) {
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

func WithBody(v any) RequestOpt {
	return func(r *Request) error {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}

		fmt.Println("### body", string(b))
		r.req.Body = io.NopCloser(bytes.NewReader(b))
		return nil
	}
}

func WithForm(form url.Values) RequestOpt {
	return func(r *Request) error {
		r.req.Body = io.NopCloser(strings.NewReader(form.Encode()))
		return nil
	}
}

func WithHeader(key, value string) RequestOpt {
	return func(r *Request) error {
		r.req.Header.Add(key, value)
		return nil
	}
}

func WithCredentials(credentials model.UserCredentials) RequestOpt {
	return func(r *Request) error {
		r.req.Header.Set("Authorization", "Bearer "+credentials.AccessToken)
		r.req.Header.Set("x-context-token", credentials.ContextToken)
		return nil
	}
}

func WithContentType(contentType string) RequestOpt {
	return func(r *Request) error {
		r.req.Header.Set("Content-Type", contentType)
		return nil
	}
}

func WithQueryParam(key, value string) RequestOpt {
	return func(r *Request) error {
		q := r.req.URL.Query()
		q.Add(key, value)
		r.req.URL.RawQuery = q.Encode()
		return nil
	}
}

type HttpError struct {
	Message string
	Code    int
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("Error %d: %s.",
		e.Code,
		e.Message,
	)
}

func newHttpError(message string, code int) *HttpError {
	return &HttpError{message, code}
}

func (r *Request) Send(v any) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(r.req)
	if err != nil {
		fmt.Println("### err", err.Error())

		return resp, err
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}

	fmt.Println("### resp", r.req.URL.String(), string(bytes))

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

	fmt.Println("### resp", req.URL.String(), string(bytes))

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

func Login(orgID, username, password string) (string, error) {
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

func UserContext(orgID, token string, acceptTerms bool) (string, error) {
	reqBody := map[string]any{
		"org_id": orgID,
	}
	if acceptTerms {
		for _, t := range []string{"accept_privacy_policy", "accept_terms_and_conditions"} {
			reqBody[t] = true
		}
	}

	url := "http://localhost:8050/v1/contexts"
	r, err := NewRequest(http.MethodPost, url, WithBody(reqBody), WithHeader("Authorization", fmt.Sprintf("Bearer %s", token)))
	if err != nil {
		return "", err
	}

	var contextResponse struct {
		Token string `json:"token"`
	}

	if _, err = r.Send(&contextResponse); err != nil {
		return "", err
	}

	return contextResponse.Token, nil
}

func orgAttributes(idpToken, contextToken string) ([]*model.OrgAttribute, error) {
	r, err := NewRequest(http.MethodGet, "http://localhost:8050/v1/attributes?page=1&itemsPerPage=10&matchStatus=ACTIVE&matchEditable=true",
		WithHeader("Authorization", fmt.Sprintf("Bearer %s", idpToken)),
		WithHeader("x-context-token", contextToken))

	if err != nil {
		return nil, err
	}

	var orgAttributesResp struct {
		Attributes []*model.OrgAttribute `json:"attributes"`
	}
	if _, err = r.Send(&orgAttributesResp); err != nil {
		return nil, err
	}

	return orgAttributesResp.Attributes, nil
}
