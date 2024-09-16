package main_suite_test

import (
	"database/sql"
	"fmt"
	"platform_e2e/pkg/client"
	"platform_e2e/pkg/client/model"
	"platform_e2e/pkg/config"
	"platform_e2e/pkg/keycloak"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type MainSuite struct {
	suite.Suite

	db         *sql.DB
	apiClient  client.APIClient
	keycloak   *keycloak.Client
	superAdmin *model.UserCredentials

	org        *model.Organization
	otherOrg   *model.Organization
	course     *model.Course
	otherAdmin *model.UserCredentials
	orgAdmin   *model.UserCredentials

	orgAdminInfo   *model.User
	learnerInfo    *model.User
	otherAdminInfo *model.User

	learningPlan  *model.LearningPlan
	learningGroup *model.LearningGroup

	quiz *model.LearningItem
}

func openDB(config *config.Config) (*sql.DB, error) {
	cfg := mysql.Config{
		User:                 config.DBUser,
		Passwd:               config.DBPassword,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", config.DBHost, config.DBPort),
		AllowNativePasswords: true,
	}
	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	return db, err
}

func (s *MainSuite) setupOrgUsers() {
	var err error
	s.orgAdminInfo, err = s.apiClient.CreateUser(&model.CreateUserRequest{
		Email:     "dborry+e2e_admin@learntowin.com",
		FirstName: "David",
		LastName:  "Borry",
		Roles:     []string{"ROLE_ADMIN"},
	}, s.superAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(s.orgAdminInfo.ID)

	s.learnerInfo, err = s.apiClient.CreateUser(&model.CreateUserRequest{
		Email:     "dborry+e2e_learner@learntowin.com",
		FirstName: "David",
		LastName:  "Borry",
		Roles:     []string{"ROLE_LEARNER"},
	}, s.superAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(s.learnerInfo.ID)
}

func (s *MainSuite) setupLearningGroup() {
	fmt.Println(".setupLearningGroup", s.orgAdmin)
	err := s.apiClient.CreateOrgAttribute(&model.CreateOrgAttributeRequest{
		AttributeOptions: []*model.AttributeOption{
			{Label: "Red", SequenceOrder: 0},
			{Label: "Green", SequenceOrder: 1},
			{Label: "Blue", SequenceOrder: 2},
		},
		Name:         "Color",
		Organization: fmt.Sprintf("/api/organizations/%s", s.org.ID),
		Type:         "SINGLE_SELECT",
	}, s.orgAdmin)
	s.Require().Nil(err)

	attributes, err := s.apiClient.OrgAttributes(s.orgAdmin)
	s.Require().Nil(err)
	s.Require().Equal(1, len(attributes))

	attributeID := attributes[0].ID
	s.Require().NotEmpty(attributeID)

	err = s.apiClient.AssignUserAttribute(&model.AssignUserAttributeRequest{
		UserID:      s.learnerInfo.ID,
		AttributeID: attributeID,
		Value:       "Blue",
	}, s.orgAdmin)
	s.Require().Nil(err)

	s.learningGroup, err = s.apiClient.CreateLearningGroup(&model.CreateLearningGroupRequest{
		Name: "Test LG",
		Attributes: []*model.AttributeFilter{
			{AttributeID: attributeID, FilterOperator: "EQ", Value: "Blue"},
		},
	}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().Equal(1, s.learningGroup.UserCount)
}

func (s *MainSuite) setupCourse() {
	var err error
	s.course, err = s.apiClient.CreateCourse(&model.CreateCourseRequest{OrganizationId: s.org.ID, VersionName: "Geography", Title: "Geography"}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(s.course.ID)

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

	cards, err := s.apiClient.CreateCardsFromFile(lesson.ID, "./testdata/cards.json", s.orgAdmin)
	s.Require().Nil(err)
	s.Require().Equal(3, len(cards))

	cards, err = s.apiClient.LearningItemCards(lesson.ID, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().Equal(3, len(cards))

	/*media, err := s.apiClient.createMedia(createMediaRequest{FileName: "image.jpg", MimeType: "image/jpeg"}, s.orgAdmin)
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
	introCard, err = s.apiClient.updateCard(updateCardRequest{
		ID:    introCard.ID,
		JSON:  introCard.JSON,
		Title: "Test",
		Media: []string{"/api/media/" + media.ID},
	}, s.orgAdmin)
	s.Require().Nil(err)*/

	s.quiz, err = s.apiClient.CreateLearningItem(&model.CreateLearningItemRequest{
		Course:      "/api/courses/" + s.course.ID,
		Type:        "quiz",
		State:       "draft",
		Name:        "Geography Test",
		Description: "It's time to test your knowledge of geography in this quiz.",
		Points:      1,
	}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(s.quiz.ID)

	cards, err = s.apiClient.CreateCardsFromFile(s.quiz.ID, "./testdata/quiz.json", s.orgAdmin)
	s.Assert().Nil(err)
	s.Assert().Equal(6, len(cards))

	s.learningPlan, err = s.apiClient.CreateLearningPlan(&model.CreateLearningPlanRequest{Name: "Semester 1", ActivatedAt: time.Now().Format(time.RFC3339)}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(s.learningPlan.ID)
	s.Require().Empty(s.learningPlan.Courses)

	s.learningPlan, err = s.apiClient.AddCoursesToLearningPlan(&model.AddCoursesToLearningPlanRequest{LearningPlanID: s.learningPlan.ID, Courses: []string{s.course.ID}}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(s.learningPlan.ID)
	s.Require().NotEmpty(s.learningPlan.Courses)

	s.learningPlan, err = s.apiClient.AddGroupsToLearningPlan(&model.AddGroupsToLearningPlanRequest{LearningPlanID: s.learningPlan.ID, LearningGroupIDs: []string{s.learningGroup.ID}}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(s.learningPlan.ID)
	s.Require().NotEmpty(s.learningPlan.LearningGroups)

	_, err = s.apiClient.ActivateCourse(s.course.ID, s.orgAdmin)
	s.Require().Nil(err)

	_, err = s.apiClient.ActivateLearningPlan(s.learningPlan.ID, s.orgAdmin)
	s.Require().Nil(err)

	cards, err = s.apiClient.LearningItemCards(lesson.ID, s.orgAdmin)
	s.Require().Nil(err)

	/*newIntroCard := cards[0]
	s.Require().Equal(introCard.ID, newIntroCard.ID)
	imgBlock := newIntroCard.JSON.ContentBlocks[2]
	s.Require().Equal("Test", newIntroCard.Title)

	s.Require().NotNil(imgBlock.MediaID)
	s.Require().Equal(media.ID, *imgBlock.MediaID)*/

}

func (s *MainSuite) SetupSuite() {
	err := godotenv.Load("../../.env")
	s.Require().Nil(err)

	cnf := config.Get()
	s.db, err = openDB(cnf)
	s.Require().Nil(err)

	s.keycloak = keycloak.New()
	s.apiClient = client.New(cnf)

	s.superAdmin, err = s.apiClient.Login(cnf.SuperAdminEmail, cnf.SuperAdminPassword, cnf.SuperAdminOrgID, false)
	s.Require().Nil(err)

	s.org, err = s.apiClient.CreateOrganization(&model.CreateOrganizationRequest{Slug: "atg", Name: "atg", Status: "ACTIVE"}, s.superAdmin)
	s.Require().Nil(err)
	s.Require().Equal(s.org.Slug, "atg")

	err = s.apiClient.SwitchOrg(s.org.ID, s.superAdmin)
	s.Require().Nil(err)

	s.setupOrgUsers()
	s.orgAdmin, err = s.apiClient.Login("dborry+e2e_admin@learntowin.com", "password", s.org.ID, true)
	s.Require().Nil(err)

	s.setupLearningGroup()
	s.setupCourse()

	// Other org setup
	s.otherOrg, err = s.apiClient.CreateOrganization(&model.CreateOrganizationRequest{Slug: "other-org", Name: "Other Org", Status: "ACTIVE"}, s.superAdmin)
	s.Require().Nil(err)
	s.Require().Equal(s.org.Slug, "atg")

	s.apiClient.SwitchOrg(s.otherOrg.ID, s.superAdmin)
	s.otherAdminInfo, err = s.apiClient.CreateUser(&model.CreateUserRequest{
		Email:     "dborry+e2e_other_learner@learntowin.com",
		FirstName: "David",
		LastName:  "Borry",
		Roles:     []string{"ROLE_ADMIN"},
	}, s.superAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(s.otherAdminInfo.ID)

	s.otherAdmin, err = s.apiClient.Login(s.otherAdminInfo.Email, "password", s.otherOrg.ID, true)
	s.Require().Nil(err)

}

func (s *MainSuite) TearDownSuite() {
	for _, o := range []*model.Organization{s.org, s.otherOrg} {
		_, err := s.db.Exec("delete from organization.organization where slug = ?", o.Slug)
		s.Assert().Nil(err)
		err = s.keycloak.DeleteRealm(o.ID)
		s.Assert().Nil(err)
	}

}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(MainSuite))
}

func ptr[T any](t T) *T {
	return &t
}
