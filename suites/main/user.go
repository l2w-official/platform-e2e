package main_suite_test

type userCredentials struct {
	accessToken  string
	contextToken string
}

func (c *userCredentials) switchOrg(orgID string) (err error) {
	if c.contextToken, err = userContext(orgID, c.accessToken, false); err != nil {
		return err
	}

	return nil
}

func userLogin(username, password, orgID string, acceptTerms bool) (c userCredentials, err error) {
	if c.accessToken, err = login(orgID, username, password); err != nil {
		return c, err
	}

	if c.contextToken, err = userContext(orgID, c.accessToken, acceptTerms); err != nil {
		return c, err
	}

	return c, nil
}
