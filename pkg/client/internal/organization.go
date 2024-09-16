package internal

import (
	"net/http"
	"platform_e2e/pkg/client/model"
)

func (cli *APIClient) CreateOrganization(req *model.CreateOrganizationRequest, credentials *model.UserCredentials) (*model.Organization, error) {
	var resp model.CreateOrganizationResponse
	if err := cli.sendRequest(http.MethodPost, "/v1/organizations", req, credentials, &resp); err != nil {
		return nil, err
	}

	return &resp.Organization, nil
}

func (cli *APIClient) CreateUser(req *model.CreateUserRequest, credentials *model.UserCredentials) (*model.User, error) {
	var user *model.User
	if err := cli.sendRequest(http.MethodPost, "/v1/users", req, credentials, &user); err != nil {
		return nil, err
	}

	return user, nil

}

func (cli *APIClient) CreateOrgAttribute(req *model.CreateOrgAttributeRequest, credentials *model.UserCredentials) error {
	return cli.sendRequest(http.MethodPost, "/v1/attributes", req, credentials, nil)
}

func (cli *APIClient) OrgAttributes(credentials *model.UserCredentials) ([]*model.OrgAttribute, error) {
	var orgAttributesResp struct {
		Attributes []*model.OrgAttribute `json:"attributes"`
	}

	if err := cli.sendRequest(http.MethodGet, "/v1/attributes?page=1&itemsPerPage=10&matchStatus=ACTIVE&matchEditable=true", nil, credentials, &orgAttributesResp); err != nil {
		return nil, err
	}

	return orgAttributesResp.Attributes, nil
}

func (cli *APIClient) AssignUserAttribute(req *model.AssignUserAttributeRequest, credentials *model.UserCredentials) error {
	return cli.sendRequest(http.MethodPost, "/v1/user_attributes", req, credentials, nil)
}
