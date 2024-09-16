package internal

import (
	"encoding/json"
	"net/http"
	"os"
	"platform_e2e/pkg/client/model"
	"platform_e2e/pkg/transport"
)

func (cli *APIClient) CreateLearningGroup(req *model.CreateLearningGroupRequest, credentials *model.UserCredentials) (*model.LearningGroup, error) {
	var lg model.LearningGroup
	if err := cli.sendRequest(http.MethodPost, "/v1/learning_groups", req, credentials, &lg); err != nil {
		return nil, err
	}

	return &lg, nil
}

func (cli *APIClient) CreateCourse(req *model.CreateCourseRequest, credentials *model.UserCredentials) (*model.Course, error) {
	r, err := transport.NewRequest(http.MethodPost, cli.URL+"/v1/courses", transport.WithBody(req), transport.WithCredentials(*credentials), transport.WithContentType("application/ld+json"))
	if err != nil {
		return nil, err
	}

	var course model.Course
	if _, err = r.Send(&course); err != nil {
		return nil, err
	}

	return &course, nil
}

func (cli *APIClient) CreateLearningItem(req *model.CreateLearningItemRequest, credentials *model.UserCredentials) (*model.LearningItem, error) {
	r, err := transport.NewRequest(http.MethodPost, cli.URL+"/v1/learning_items", transport.WithBody(req), transport.WithCredentials(*credentials), transport.WithContentType("application/ld+json"))
	if err != nil {
		return nil, err
	}

	var learningItem model.LearningItem
	if _, err = r.Send(&learningItem); err != nil {
		return nil, err
	}

	return &learningItem, nil
}

func (cli *APIClient) CreateCard(req *model.CreateCardRequest, credentials *model.UserCredentials) (*model.Card, error) {
	r, err := transport.NewRequest(http.MethodPost, cli.URL+"/v1/cards", transport.WithBody(req), transport.WithCredentials(*credentials), transport.WithContentType("application/ld+json"))
	if err != nil {
		return nil, err
	}

	var card model.Card
	if _, err = r.Send(&card); err != nil {
		return nil, err
	}

	return &card, nil
}

func (cli *APIClient) CreateCardsFromFile(learningItemID string, cardsFilepath string, credentials *model.UserCredentials) ([]*model.Card, error) {
	b, err := os.ReadFile(cardsFilepath)
	if err != nil {
		return nil, err
	}

	var req model.CreateCardsRequest
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	var resp struct {
		Cards []*model.Card `json:"cards"`
	}

	if err = cli.sendRequest(http.MethodPost, "/v1/learning_items/"+learningItemID+"/cards", req, credentials, &resp); err != nil {
		return nil, err
	}

	return resp.Cards, nil
}

func (cli *APIClient) LearningItemCards(learningItemID string, credentials *model.UserCredentials) ([]*model.Card, error) {
	var resp struct {
		Cards []*model.Card `json:"hydra:member"`
	}

	if err := cli.sendRequest(http.MethodGet, "/v1/learning_items/"+learningItemID+"/cards?order[sequenceOrder]=asc", nil, credentials, &resp); err != nil {
		return nil, err
	}

	return resp.Cards, nil
}

func (cli *APIClient) CreateLearningPlan(req *model.CreateLearningPlanRequest, credentials *model.UserCredentials) (*model.LearningPlan, error) {
	var learningPlan model.LearningPlan
	if err := cli.sendRequest(http.MethodPost, "/v1/learning_plans", req, credentials, &learningPlan); err != nil {
		return nil, err
	}

	return &learningPlan, nil
}

func (cli *APIClient) AddCoursesToLearningPlan(req *model.AddCoursesToLearningPlanRequest, credentials *model.UserCredentials) (*model.LearningPlan, error) {
	var learningPlan model.LearningPlan
	if err := cli.sendRequest(http.MethodPost, "/v1/learning_plan/"+req.LearningPlanID+"/courses", req, credentials, &learningPlan); err != nil {
		return nil, err
	}

	return &learningPlan, nil
}

func (cli *APIClient) AddGroupsToLearningPlan(req *model.AddGroupsToLearningPlanRequest, credentials *model.UserCredentials) (*model.LearningPlan, error) {
	var learningPlan model.LearningPlan
	if err := cli.sendRequest(http.MethodPost, "/v1/learning_plans/"+req.LearningPlanID+"/learning_plan_groups", req, credentials, &learningPlan); err != nil {
		return nil, err
	}

	return &learningPlan, nil
}

func (cli *APIClient) ActivateCourse(courseID string, credentials *model.UserCredentials) (*model.Course, error) {
	var course model.Course
	if err := cli.sendRequest(http.MethodPatch, "/v1/courses/"+courseID, map[string]any{
		"state": "published",
	}, credentials, &course, transport.WithContentType("application/merge-patch+json")); err != nil {
		return nil, err
	}

	return &course, nil
}

func (cli *APIClient) ActivateLearningPlan(learningPlanID string, credentials *model.UserCredentials) (*model.LearningPlan, error) {
	var learningPlan model.LearningPlan
	if err := cli.sendRequest(http.MethodPatch, "/v1/learning_plans/"+learningPlanID, map[string]any{
		"state": 1,
	}, credentials, &learningPlan, transport.WithContentType("application/merge-patch+json")); err != nil {
		return nil, err
	}

	return &learningPlan, nil
}

func (cli *APIClient) CreateMedia(req *model.CreateMediaRequest, credentials *model.UserCredentials) (*model.Media, error) {
	var resp *model.Media
	if err := cli.sendRequest(http.MethodPost, "/v1/media", req, credentials, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

func (cli *APIClient) UpdateCard(req *model.UpdateCardRequest, credentials *model.UserCredentials) (*model.Card, error) {
	var resp model.Card
	if err := cli.sendRequest(http.MethodPatch, "/v1/cards/"+req.ID, req, credentials, &resp, transport.WithContentType("application/merge-patch+json")); err != nil {
		return nil, err
	}

	return &resp, nil
}
