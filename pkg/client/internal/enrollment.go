package internal

import (
	"net/http"
	"platform_e2e/pkg/client/model"
	"platform_e2e/pkg/transport"
)

func (cli *APIClient) Invitations(filter *model.InvitationsFilter, credentials *model.UserCredentials) ([]*model.Invitation, error) {
	var resp struct {
		Invitations []*model.Invitation `json:"hydra:member"`
	}

	if err := cli.sendRequest(http.MethodGet, "/v1/invitations", nil, credentials, &resp,
		transport.WithQueryParam("courseId[]", filter.CourseID),
		transport.WithQueryParam("invitedUserId[]", filter.InvitedUserID)); err != nil {
		return nil, err
	}

	return resp.Invitations, nil
}
