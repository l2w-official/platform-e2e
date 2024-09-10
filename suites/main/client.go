package main_suite_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type apiClient struct {
	url string
}

func newApiClient() *apiClient {
	return &apiClient{
		url: config.apiGatewayURL,
	}
}

func (cli *apiClient) sendRequest(method, path string, req any, credentials userCredentials, res any, opts ...requestOpt) error {
	opts = append([]requestOpt{withBody(req), withCredentials(credentials), withContentType("application/ld+json")}, opts...)
	r, err := newRequest(method, cli.url+path, opts...)
	if err != nil {
		return err
	}

	if _, err = r.send(res); err != nil {
		return err
	}

	return nil
}

type createOrganizationRequest struct {
	Slug   string `json:"slug"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type createOrganizationResponse struct {
	Organization organization `json:"organization"`
}

func (cli *apiClient) createOrganization(req createOrganizationRequest, credentials userCredentials) (*organization, error) {
	var resp createOrganizationResponse
	if err := cli.sendRequest(http.MethodPost, "/v1/organizations", req, credentials, &resp); err != nil {
		return nil, err
	}

	return &resp.Organization, nil
}

type user struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Roles     []string `json:"roles"`
}

type createUserRequest struct {
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Roles     []string `json:"roles"`
}

func (cli *apiClient) createUser(req createUserRequest, credentials userCredentials) (*user, error) {
	var user *user
	if err := cli.sendRequest(http.MethodPost, "/v1/users", req, credentials, &user); err != nil {
		return nil, err
	}

	return user, nil

}

type attributeOption struct {
	Label         string `json:"label"`
	SequenceOrder int    `json:"sequenceOrder"`
}
type createOrgAttributeRequest struct {
	AttributeOptions []*attributeOption `json:"attributeOptions"`
	Name             string             `json:"name"`
	Organization     string             `json:"organization"`
	Type             string             `json:"type"`
}

func (cli *apiClient) createOrgAttribute(req createOrgAttributeRequest, credentials userCredentials) error {
	return cli.sendRequest(http.MethodPost, "/v1/attributes", req, credentials, nil)
}

func (cli *apiClient) orgAttributes(credentials userCredentials) ([]*orgAttribute, error) {
	var orgAttributesResp struct {
		Attributes []*orgAttribute `json:"attributes"`
	}

	if err := cli.sendRequest(http.MethodGet, "/v1/attributes?page=1&itemsPerPage=10&matchStatus=ACTIVE&matchEditable=true", nil, credentials, &orgAttributesResp); err != nil {
		return nil, err
	}

	return orgAttributesResp.Attributes, nil
}

type assignUserAttributesRequest struct {
	UserID      string `json:"userId"`
	AttributeID string `json:"attributeId"`
	Value       string `json:"value"`
}

func (cli *apiClient) assignUserAttributes(req assignUserAttributesRequest, credentials userCredentials) error {
	return cli.sendRequest(http.MethodPost, "/v1/user_attributes", req, credentials, nil)
}

type attributeFilter struct {
	AttributeID    string `json:"attributeId"`
	FilterOperator string `json:"filterOperator"`
	Value          string `json:"value"`
}
type createLearningGroupRequest struct {
	Name       string             `json:"name"`
	Attributes []*attributeFilter `json:"attributes"`
}

type learningGroup struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	UserCount int    `json:"userCount"`
}

func (cli *apiClient) createLearningGroup(req createLearningGroupRequest, credentials userCredentials) (*learningGroup, error) {
	var lg learningGroup
	if err := cli.sendRequest(http.MethodPost, "/v1/learning_groups", req, credentials, &lg); err != nil {
		return nil, err
	}

	return &lg, nil
}

type createCourseRequest struct {
	OrganizationId string `json:"organizationId"`
	Title          string `json:"title"`
	VersionName    string `json:"versionName"`
}

type course struct {
	ID string `json:"id"`
}

func (cli *apiClient) createCourse(req createCourseRequest, credentials userCredentials) (*course, error) {
	r, err := newRequest(http.MethodPost, cli.url+"/v1/courses", withBody(req), withCredentials(credentials), withContentType("application/ld+json"))
	if err != nil {
		return nil, err
	}

	var course course
	if _, err = r.send(&course); err != nil {
		return nil, err
	}

	return &course, nil
}

type createLearningItemRequest struct {
	Course        string `json:"course"`
	Description   string `json:"description"`
	Name          string `json:"name"`
	Points        int    `json:"points"`
	SequenceOrder int    `json:"sequenceOrder"`
	State         string `json:"state"`
	Type          string `json:"type"`
}

type learningItem struct {
	ID                    string `json:"id"`
	LearningItemVersionID string `json:"learningItemVersionId"`
}

func (cli *apiClient) createLearningItem(req createLearningItemRequest, credentials userCredentials) (*learningItem, error) {
	r, err := newRequest(http.MethodPost, cli.url+"/v1/learning_items", withBody(req), withCredentials(credentials), withContentType("application/ld+json"))
	if err != nil {
		return nil, err
	}

	var learningItem learningItem
	if _, err = r.send(&learningItem); err != nil {
		return nil, err
	}

	return &learningItem, nil
}

type cardContentBlock struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	JSON any    `json:"json,omitempty"`

	MediaID *string `json:"mediaId,omitempty"`
	Name    *string `json:"name,omitempty"`
}
type cardJSON struct {
	Version       string             `json:"version"`
	Description   string             `json:"description"`
	TemplateType  *string            `json:"templateType"`
	ContentBlocks []cardContentBlock `json:"contentBlocks"`
}
type createCardRequest struct {
	LearningItem    string `json:"learningItem,omitempty"`
	Type            string `json:"type"`
	Title           string `json:"title"`
	SequenceOrder   int    `json:"sequenceOrder"`
	ConfidenceCheck bool   `json:"confidenceCheck"`

	JSON any `json:"json"`
}

type card struct {
	ID   string   `json:"id"`
	JSON cardJSON `json:"json"`
}

var i = 0

func (cli *apiClient) createCard(req createCardRequest, credentials userCredentials) (*card, error) {
	b, _ := json.Marshal(req)
	os.WriteFile(fmt.Sprintf("./testdata/card_%d.json", i), b, 0644)
	i += 1
	r, err := newRequest(http.MethodPost, cli.url+"/v1/cards", withBody(req), withCredentials(credentials), withContentType("application/ld+json"))
	if err != nil {
		return nil, err
	}

	var card card
	if _, err = r.send(&card); err != nil {
		return nil, err
	}

	return &card, nil
}

type createCardsRequest struct {
	Cards []*createCardRequest `json:"cards"`
}

func (cli *apiClient) createCardsFromFile(learningItemID string, cardsFilepath string, credentials userCredentials) ([]*card, error) {
	b, err := os.ReadFile(cardsFilepath)
	if err != nil {
		return nil, err
	}

	var req createCardsRequest
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	var resp struct {
		Cards []*card `json:"cards"`
	}

	if err = cli.sendRequest(http.MethodPost, "/v1/learning_items/"+learningItemID+"/cards", req, credentials, &resp); err != nil {
		return nil, err
	}

	return resp.Cards, nil
}

type createLearningPlanRequest struct {
	Name        string `json:"name"`
	ActivatedAt string `json:"activatedAt"`
}

type learningPlan struct {
	ID             string          `json:"id"`
	Courses        []course        `json:"courses"`
	LearningGroups []learningGroup `json:"learningGroups"`
}

func (cli *apiClient) createLearningPlan(req createLearningPlanRequest, credentials userCredentials) (*learningPlan, error) {
	var learningPlan learningPlan
	if err := cli.sendRequest(http.MethodPost, "/v1/learning_plans", req, credentials, &learningPlan); err != nil {
		return nil, err
	}

	return &learningPlan, nil
}

type addCoursesToLearningPlanRequest struct {
	Courses []string `json:"courses"`
}

func (cli *apiClient) addCoursesToLearningPlan(learningPlanID string, req addCoursesToLearningPlanRequest, credentials userCredentials) (*learningPlan, error) {
	var learningPlan learningPlan
	if err := cli.sendRequest(http.MethodPost, "/v1/learning_plan/"+learningPlanID+"/courses", req, credentials, &learningPlan); err != nil {
		return nil, err
	}

	return &learningPlan, nil
}

type addGroupsToLearningPlanRequest struct {
	LearningGroupIDs []string `json:"learningGroupIds"`
}

func (cli *apiClient) addGroupsToLearningPlan(learningPlanID string, req addGroupsToLearningPlanRequest, credentials userCredentials) (*learningPlan, error) {
	var learningPlan learningPlan
	if err := cli.sendRequest(http.MethodPost, "/v1/learning_plans/"+learningPlanID+"/learning_plan_groups", req, credentials, &learningPlan); err != nil {
		return nil, err
	}

	return &learningPlan, nil
}

func (cli *apiClient) activateCourse(courseID string, credentials userCredentials) (*course, error) {
	var course course
	if err := cli.sendRequest(http.MethodPatch, "/v1/courses/"+courseID, map[string]any{
		"state": "published",
	}, credentials, &course, withContentType("application/merge-patch+json")); err != nil {
		return nil, err
	}

	return &course, nil
}

func (cli *apiClient) activateLearningPlan(learningPlanID string, credentials userCredentials) (*learningPlan, error) {
	var learningPlan learningPlan
	if err := cli.sendRequest(http.MethodPatch, "/v1/learning_plans/"+learningPlanID, map[string]any{
		"state": 1,
	}, credentials, &learningPlan, withContentType("application/merge-patch+json")); err != nil {
		return nil, err
	}

	return &learningPlan, nil
}

// Offline
type courseBundleRequest struct {
	OrgID          string `json:"orgId"`
	CourseID       string `json:"courseId"`
	LearningPlanID string `json:"learningplanId"`
	DeviceID       string `json:"deviceId"`
	OfflineMode    string `json:"offlineMode"`
}

type courseBundleResponse struct {
	JobID string `json:"jobId"`
}

func (cli *apiClient) courseBundle(req courseBundleRequest, credentials userCredentials) (*courseBundleResponse, error) {
	var resp courseBundleResponse
	if err := cli.sendRequest(http.MethodPost, "/v1/course-bundle", req, credentials, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type courseBundleURLResponse struct {
	CourseURL    string `json:"courseUrl"`
	BundleStatus string `json:"bundleStatus"`
}

func (cli *apiClient) courseBundleURL(jobID string, credentials userCredentials) (*courseBundleURLResponse, error) {
	var resp courseBundleURLResponse
	if err := cli.sendRequest(http.MethodGet, "/v1/course-bundle-url/"+jobID, nil, credentials, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type invitation struct {
	ID                string `json:"id"`
	InvitedUserID     string `json:"invitedUserId"`
	DownloadedOffline bool   `json:"downloadedOffline"`
}

type invitationsFilter struct {
	courseID      string
	invitedUserID string
}

func (cli *apiClient) invitations(filter invitationsFilter, credentials userCredentials) ([]*invitation, error) {
	var resp struct {
		Invitations []*invitation `json:"hydra:member"`
	}

	if err := cli.sendRequest(http.MethodGet, "/v1/invitations", nil, credentials, &resp,
		withQueryParam("courseId[]", filter.courseID),
		withQueryParam("invitedUserId[]", filter.invitedUserID)); err != nil {
		return nil, err
	}

	return resp.Invitations, nil
}

type invitationEnrollRequest struct {
	InvitationID string `json:"invitationId"`
}
type invitationEnrollResponse struct {
	ID           string `json:"id"`
	InvitationID string `json:"invitationId"`
	Status       string `json:"status"`
}

func (cli *apiClient) invitationEnroll(req invitationEnrollRequest, credentials userCredentials) (*invitationEnrollResponse, error) {
	var resp invitationEnrollResponse
	if err := cli.sendRequest(http.MethodPost, "/v1/invitation-enroll", req, credentials, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (cli *apiClient) enrollmentJob(jobID string, credentials userCredentials) (*invitationEnrollResponse, error) {
	var resp invitationEnrollResponse
	if err := cli.sendRequest(http.MethodGet, "/v1/invitation-enroll/"+jobID, nil, credentials, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type cloneEnrollmentRequest struct {
	InvitationID string `json:"invitationId"`
}

type cardEnrollment struct {
	CardEnrollmentId         string   `json:"cardEnrollmentId,omitempty"`
	LearningItemEnrollmentId string   `json:"learningItemEnrollmentId,omitempty"`
	CardId                   string   `json:"cardId,omitempty"`
	Score                    int32    `json:"score,omitempty"`
	ElapsedSec               int32    `json:"elapsedSec,omitempty"`
	ServerEnrollment         bool     `json:"serverEnrollment,omitempty"`
	DeviceID                 string   `json:"deviceId,omitempty"`
	Approved                 bool     `json:"approved,omitempty"`
	Answer                   []string `json:"answer,omitempty"`
	Confidence               int32    `json:"confidence,omitempty"`
	CreatedAt                string   `json:"createdAt,omitempty"`
	UpdatedAt                string   `json:"updatedAt,omitempty"`
	StartedAt                string   `json:"startedAt,omitempty"`
	CompletedAt              string   `json:"completedAt,omitempty"`
	Progress                 int32    `json:"progress,omitempty"`
	TotalPoints              int32    `json:"totalPoints,omitempty"`
}
type learningItemEnrollment struct {
	LearningItemEnrollmentId string            `json:"learningItemEnrollmentId,omitempty"`
	CourseEnrollmentId       string            `json:"courseEnrollmentId,omitempty"`
	LearningItemId           string            `json:"learningItemId,omitempty"`
	DeviceId                 string            `json:"deviceId,omitempty"`
	StartedAt                string            `json:"startedAt,omitempty"`
	UpdatedAt                string            `json:"updatedAt,omitempty"`
	Progress                 int32             `json:"progress,omitempty"`
	CardEnrollments          []*cardEnrollment `json:"cardEnrollments,omitempty"`
}
type cloneEnrollmentResponse struct {
	CourseEnrollmentID      string                    `json:"courseEnrollmentId"`
	CourseID                string                    `json:"courseId"`
	InvitationID            string                    `json:"invitationId"`
	LearningItemEnrollments []*learningItemEnrollment `json:"learningItemEnrollments"`
	UserID                  string                    `json:"userId"`
}

func (cli *apiClient) cloneEnrollment(req cloneEnrollmentRequest, credentials userCredentials) (*cloneEnrollmentResponse, error) {
	var resp cloneEnrollmentResponse
	if err := cli.sendRequest(http.MethodPost, "/v1/enrollments/clone", req, credentials, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type syncEnrollmentRequest struct {
	LearningItemEnrollments []*learningItemEnrollment `json:"LearningItemEnrollments"`
}

type syncEnrollmentsResult struct {
	LearningItemEnrollmentID string `json:"learningItemEnrollmentId"`
	Success                  bool   `json:"success"`
	Message                  string `json:"message"`
}

func (cli *apiClient) syncEnrollments(req syncEnrollmentRequest, credentials userCredentials) ([]*syncEnrollmentsResult, error) {
	var resp struct {
		Results []*syncEnrollmentsResult `json:"results"`
	}
	if err := cli.sendRequest(http.MethodPost, "/v1/enrollments/sync", req, credentials, &resp); err != nil {
		return nil, err
	}

	return resp.Results, nil
}

type getDuplicateEnrollmentsRequest struct {
	LearningItemEnrollmentIDs []string `json:"learningItemEnrollmentIds"`
}

func (cli *apiClient) getDuplicateEnrollments(req getDuplicateEnrollmentsRequest, credentials userCredentials) ([]*learningItemEnrollment, error) {
	var resp struct {
		LearningItemEnrollments []*learningItemEnrollment `json:"learningItemEnrollments"`
	}
	if err := cli.sendRequest(http.MethodPost, "/v1/enrollments/duplicate", req, credentials, &resp); err != nil {
		return nil, err
	}

	return resp.LearningItemEnrollments, nil
}
