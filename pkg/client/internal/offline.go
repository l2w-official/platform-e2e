package internal

import (
	"net/http"
	"platform_e2e/pkg/client/model"
	"platform_e2e/pkg/transport"
)

func (cli *APIClient) CourseBundle(req *model.CourseBundleRequest, credentials *model.UserCredentials) (*model.CourseBundleResponse, error) {
	var resp model.CourseBundleResponse
	if err := cli.sendRequest(http.MethodPost, "/v1/course-bundle", req, credentials, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (cli *APIClient) CourseBundleURL(jobID string, credentials *model.UserCredentials) (*model.CourseBundleURLResponse, error) {
	var resp model.CourseBundleURLResponse
	if err := cli.sendRequest(http.MethodGet, "/v1/course-bundle-url/"+jobID, nil, credentials, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (cli *APIClient) InvitationEnroll(req *model.InvitationEnrollRequest, credentials *model.UserCredentials) (*model.InvitationEnrollResponse, error) {
	var resp model.InvitationEnrollResponse
	if err := cli.sendRequest(http.MethodPost, "/v1/invitation-enroll", req, credentials, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (cli *APIClient) EnrollmentJob(jobID string, credentials *model.UserCredentials) (*model.InvitationEnrollResponse, error) {
	var resp model.InvitationEnrollResponse
	if err := cli.sendRequest(http.MethodGet, "/v1/invitation-enroll/"+jobID, nil, credentials, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (cli *APIClient) CloneEnrollment(req *model.CloneEnrollmentRequest, credentials *model.UserCredentials) (*model.CloneEnrollmentResponse, error) {
	var resp model.CloneEnrollmentResponse
	if err := cli.sendRequest(http.MethodPost, "/v1/enrollments/clone", req, credentials, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (cli *APIClient) SyncEnrollments(req *model.SyncEnrollmentRequest, credentials *model.UserCredentials) ([]*model.SyncEnrollmentsResult, error) {
	var resp struct {
		Results []*model.SyncEnrollmentsResult `json:"results"`
	}
	if err := cli.sendRequest(http.MethodPost, "/v1/enrollments/sync", req, credentials, &resp); err != nil {
		return nil, err
	}

	return resp.Results, nil
}

func (cli *APIClient) DuplicateEnrollments(req *model.DuplicateEnrollmentsRequest, credentials *model.UserCredentials) ([]*model.LearningItemEnrollment, error) {
	var resp struct {
		LearningItemEnrollments []*model.LearningItemEnrollment `json:"learningItemEnrollments"`
	}
	if err := cli.sendRequest(http.MethodPost, "/v1/enrollments/duplicate", req, credentials, &resp); err != nil {
		return nil, err
	}

	return resp.LearningItemEnrollments, nil
}

func (cli *APIClient) ApproveEnrollment(req *model.ApproveEnrollmentRequest, credentials *model.UserCredentials) (bool, error) {
	var resp struct {
		Success bool `json:"success"`
	}
	if err := cli.sendRequest(http.MethodPatch, "/v1/enrollment/approve", req, credentials, &resp, transport.WithContentType("application/merge-patch+json")); err != nil {
		return false, err
	}

	return resp.Success, nil
}

func (cli *APIClient) SubmitEnrollments(req *model.SubmitEnrollmentsRequest, credentials *model.UserCredentials) ([]*model.EnrollmentStatus, error) {
	var resp struct {
		EnrollmentStatuses []*model.EnrollmentStatus `json:"enrollmentStatuses"`
	}
	if err := cli.sendRequest(http.MethodPatch, "/v1/submit-enrollments", req, credentials, &resp); err != nil {
		return nil, err
	}

	return resp.EnrollmentStatuses, nil
}
