package internal

import (
	"platform_e2e/pkg/client/model"
	"platform_e2e/pkg/transport"
)

// Login - Password login to a specific organization
func (c *APIClient) Login(username, password, orgID string, acceptTerms bool) (*model.UserCredentials, error) {
	var credentials model.UserCredentials
	var err error
	if credentials.AccessToken, err = transport.Login(orgID, username, password); err != nil {
		return nil, err
	}

	if credentials.ContextToken, err = transport.UserContext(orgID, credentials.AccessToken, acceptTerms); err != nil {
		return nil, err
	}

	return &credentials, nil
}

// SwitchOrg - Update a user session to switch to a different org
func (c *APIClient) SwitchOrg(orgID string, credentials *model.UserCredentials) (err error) {
	if credentials.ContextToken, err = transport.UserContext(orgID, credentials.AccessToken, false); err != nil {
		return err
	}

	return nil
}
