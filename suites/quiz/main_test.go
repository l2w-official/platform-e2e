package quiz_test

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"platform_e2e/pkg/client/model"
	"strings"
	"time"
)

func (s *MainSuite) uploadImage(name, filename string) (id string, err error) {
	media, err := s.apiClient.CreateMedia(&model.CreateMediaRequest{FileName: name, MimeType: "image/jpeg"}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(media.TemporaryPutURL)

	img, err := os.ReadFile(filename)
	s.Require().Nil(err)

	imgReq, err := http.NewRequest(http.MethodPut, media.TemporaryPutURL, bytes.NewReader(img))
	s.Require().Nil(err)

	imgResp, err := http.DefaultClient.Do(imgReq)
	s.Require().Nil(err)
	s.Require().Equal(imgResp.StatusCode, http.StatusOK)

	return media.ID, nil
}

func makeCardMedia(mediaIds ...string) []string {
	urls := make([]string, len(mediaIds))
	for i, mediaId := range mediaIds {
		urls[i] = "/api/media/" + mediaId
	}

	return urls
}

func (s *MainSuite) TestQuiz() {

	media1, err := s.uploadImage("france", "./testdata/fr.png")
	s.Require().NoError(err)
	media2, err := s.uploadImage("gernamy", "./testdata/de.png")
	s.Require().NoError(err)
	media3, err := s.uploadImage("us", "./testdata/us.png")
	s.Require().NoError(err)
	media4, err := s.uploadImage("brazil", "./testdata/bz.png")
	s.Require().NoError(err)

	b, err := os.ReadFile("./testdata/quiz.json")
	s.Require().NoError(err)

	cardJSON := strings.ReplaceAll(string(b), "{MEDIA_1}", media1)
	cardJSON = strings.ReplaceAll(cardJSON, "{MEDIA_2}", media2)
	cardJSON = strings.ReplaceAll(cardJSON, "{MEDIA_3}", media3)
	cardJSON = strings.ReplaceAll(cardJSON, "{MEDIA_4}", media4)

	err = os.WriteFile("./testdata/quiz_img.json", []byte(cardJSON), 0644)
	s.Require().NoError(err)

	cards, err := s.apiClient.CreateCardsFromFile(s.quiz.ID, "./testdata/quiz_img.json", s.orgAdmin)
	s.Require().Nil(err)
	s.Assert().Equal(8, len(cards))

	cardMedia := makeCardMedia(media1, media2, media3, media4)
	for _, c := range cards {
		_, err = s.apiClient.UpdateCard(&model.UpdateCardRequest{
			ID:    c.ID,
			Title: "Test",
			Media: cardMedia,
		}, s.orgAdmin)
		s.Require().Nil(err)
	}

	cards, err = s.apiClient.LearningItemCards(s.quiz.ID, s.orgAdmin)
	s.Require().NoError(err)
	for _, c := range cards {
		for _, m := range c.Media {
			fmt.Println("### MEDIA QUIZ", m.TemporaryGetURL)
		}
	}

	lesson, err := s.apiClient.CreateLearningItem(&model.CreateLearningItemRequest{
		Course:      "/api/courses/" + s.course.ID,
		Type:        "lesson",
		State:       "draft",
		Name:        "Introduction to Geography",
		Description: "Discover the basics of geography, with a fun list of famous countries and continents.",
		Points:      1,
	}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(lesson.ID)

	b, err = os.ReadFile("./testdata/cards.json")
	s.Require().NoError(err)

	cardJSON = strings.ReplaceAll(string(b), "{MEDIA_1}", media1)
	cardJSON = strings.ReplaceAll(cardJSON, "{MEDIA_2}", media2)

	err = os.WriteFile("./testdata/cards_img.json", []byte(cardJSON), 0644)
	s.Require().NoError(err)

	cards, err = s.apiClient.CreateCardsFromFile(lesson.ID, "./testdata/cards_img.json", s.orgAdmin)
	s.Require().Nil(err)
	s.Require().Equal(5, len(cards))

	cards, err = s.apiClient.LearningItemCards(lesson.ID, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().Equal(5, len(cards))

	media, err := s.apiClient.CreateMedia(&model.CreateMediaRequest{FileName: "image.jpg", MimeType: "image/jpeg"}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(media.TemporaryPutURL)

	img, err := os.ReadFile("./testdata/image.jpg")
	s.Require().Nil(err)

	imgReq, err := http.NewRequest(http.MethodPut, media.TemporaryPutURL, bytes.NewReader(img))
	s.Require().Nil(err)

	imgResp, err := http.DefaultClient.Do(imgReq)
	s.Require().Nil(err)
	s.Require().Equal(imgResp.StatusCode, http.StatusOK)

	introCard := cards[0]
	introCard.JSON.ContentBlocks[2].MediaID = &media.ID
	introCard.JSON.ContentBlocks[2].Name = ptr("Test image")
	_, err = s.apiClient.UpdateCard(&model.UpdateCardRequest{
		ID:    introCard.ID,
		JSON:  introCard.JSON,
		Title: "Test",
		Media: []string{"/api/media/" + media.ID},
	}, s.orgAdmin)
	s.Require().Nil(err)

	for _, c := range cards {
		_, err = s.apiClient.UpdateCard(&model.UpdateCardRequest{
			ID:    c.ID,
			Title: "Test",
			Media: cardMedia,
		}, s.orgAdmin)
		s.Require().Nil(err)
	}

	_, err = s.apiClient.ActivateCourse(s.course.ID, s.orgAdmin)
	s.Require().Nil(err)

	lp, err := s.apiClient.CreateLearningPlan(&model.CreateLearningPlanRequest{Name: "Semester 1", ActivatedAt: time.Now().Format(time.RFC3339)}, s.orgAdmin)
	s.Require().NoError(err)

	lp, err = s.apiClient.AddGroupsToLearningPlan(&model.AddGroupsToLearningPlanRequest{LearningPlanID: lp.ID, LearningGroupIDs: []string{s.learningGroup.ID}}, s.orgAdmin)
	s.Require().NoError(err)

	lp, err = s.apiClient.AddCoursesToLearningPlan(&model.AddCoursesToLearningPlanRequest{LearningPlanID: lp.ID, Courses: []string{s.course.ID}}, s.orgAdmin)
	s.Require().NoError(err)

	s.apiClient.ActivateLearningPlan(lp.ID, s.orgAdmin)
}
