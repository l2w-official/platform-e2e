package main_suite_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type MainSuite struct {
	suite.Suite

	db         *sql.DB
	apiClient  *apiClient
	keycloak   *keycloakCli
	superAdmin userCredentials

	org        *organization
	otherOrg   *organization
	course     *course
	otherAdmin userCredentials
	orgAdmin   userCredentials

	orgAdminInfo   *user
	learnerInfo    *user
	otherAdminInfo *user

	learningPlan  *learningPlan
	learningGroup *learningGroup

	quiz *learningItem
}

func openDB() (*sql.DB, error) {
	cfg := mysql.Config{
		User:                 config.dbUser,
		Passwd:               config.dbPassword,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", config.dbHost, config.dbPort),
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
	s.orgAdminInfo, err = s.apiClient.createUser(createUserRequest{
		Email:     "dborry+e2e_admin@learntowin.com",
		FirstName: "David",
		LastName:  "Borry",
		Roles:     []string{"ROLE_ADMIN"},
	}, s.superAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(s.orgAdminInfo.ID)

	s.learnerInfo, err = s.apiClient.createUser(createUserRequest{
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
	err := s.apiClient.createOrgAttribute(createOrgAttributeRequest{
		AttributeOptions: []*attributeOption{
			{Label: "Red", SequenceOrder: 0},
			{Label: "Green", SequenceOrder: 1},
			{Label: "Blue", SequenceOrder: 2},
		},
		Name:         "Color",
		Organization: fmt.Sprintf("/api/organizations/%s", s.org.ID),
		Type:         "SINGLE_SELECT",
	}, s.orgAdmin)
	s.Require().Nil(err)

	attributes, err := s.apiClient.orgAttributes(s.orgAdmin)
	s.Require().Nil(err)
	s.Require().Equal(1, len(attributes))

	attributeID := attributes[0].ID
	s.Require().NotEmpty(attributeID)

	err = s.apiClient.assignUserAttributes(assignUserAttributesRequest{
		UserID:      s.learnerInfo.ID,
		AttributeID: attributeID,
		Value:       "Blue",
	}, s.orgAdmin)
	s.Require().Nil(err)

	s.learningGroup, err = s.apiClient.createLearningGroup(createLearningGroupRequest{
		Name: "Test LG",
		Attributes: []*attributeFilter{
			{AttributeID: attributeID, FilterOperator: "EQ", Value: "Blue"},
		},
	}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().Equal(1, s.learningGroup.UserCount)
}

func (s *MainSuite) setupCourse() {
	var err error
	s.course, err = s.apiClient.createCourse(createCourseRequest{OrganizationId: s.org.ID, VersionName: "Geography", Title: "Geography"}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(s.course.ID)

	lesson, err := s.apiClient.createLearningItem(createLearningItemRequest{
		Course:      "/api/courses/" + s.course.ID,
		Type:        "lesson",
		State:       "draft",
		Name:        "Introduction to Geography",
		Description: "Discover the basics of geography, with a fun list of famous countries and continents.",
		Points:      1,
	}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(lesson.ID)

	cards, err := s.apiClient.createCardsFromFile(lesson.ID, "./testdata/cards.json", s.orgAdmin)
	s.Require().Nil(err)
	s.Require().Equal(3, len(cards))

	s.quiz, err = s.apiClient.createLearningItem(createLearningItemRequest{
		Course:      "/api/courses/" + s.course.ID,
		Type:        "quiz",
		State:       "draft",
		Name:        "Geography Test",
		Description: "It's time to test your knowledge of geography in this quiz.",
		Points:      1,
	}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(s.quiz.ID)

	cards, err = s.apiClient.createCardsFromFile(s.quiz.ID, "./testdata/quiz.json", s.orgAdmin)
	s.Assert().Nil(err)
	s.Assert().Equal(6, len(cards))

	s.learningPlan, err = s.apiClient.createLearningPlan(createLearningPlanRequest{Name: "Semester 1", ActivatedAt: time.Now().Format(time.RFC3339)}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(s.learningPlan.ID)
	s.Require().Empty(s.learningPlan.Courses)

	s.learningPlan, err = s.apiClient.addCoursesToLearningPlan(s.learningPlan.ID, addCoursesToLearningPlanRequest{Courses: []string{s.course.ID}}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(s.learningPlan.ID)
	s.Require().NotEmpty(s.learningPlan.Courses)

	s.learningPlan, err = s.apiClient.addGroupsToLearningPlan(s.learningPlan.ID, addGroupsToLearningPlanRequest{LearningGroupIDs: []string{s.learningGroup.ID}}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(s.learningPlan.ID)
	s.Require().NotEmpty(s.learningPlan.LearningGroups)

	_, err = s.apiClient.activateCourse(s.course.ID, s.orgAdmin)
	s.Require().Nil(err)

	_, err = s.apiClient.activateLearningPlan(s.learningPlan.ID, s.orgAdmin)
	s.Require().Nil(err)
}

func (s *MainSuite) SetupSuite() {
	err := godotenv.Load("../../.env")
	s.Require().Nil(err)

	loadConfig()
	s.db, err = openDB()
	s.Require().Nil(err)

	s.keycloak = newKeycloakCli()
	s.apiClient = newApiClient()

	s.superAdmin, err = userLogin(config.superAdminEmail, config.superAdminPassword, config.superAdminOrgID, false)
	s.Require().Nil(err)

	s.org, err = s.apiClient.createOrganization(createOrganizationRequest{Slug: "atg", Name: "atg", Status: "ACTIVE"}, s.superAdmin)
	s.Require().Nil(err)
	s.Require().Equal(s.org.Slug, "atg")

	err = s.superAdmin.switchOrg(s.org.ID)
	s.Require().Nil(err)

	s.setupOrgUsers()
	s.orgAdmin, err = userLogin("dborry+e2e_admin@learntowin.com", "password", s.org.ID, true)
	s.Require().Nil(err)

	s.setupLearningGroup()
	s.setupCourse()

	// Other org setup
	s.otherOrg, err = s.apiClient.createOrganization(createOrganizationRequest{Slug: "other-org", Name: "Other Org", Status: "ACTIVE"}, s.superAdmin)
	s.Require().Nil(err)
	s.Require().Equal(s.org.Slug, "atg")

	s.superAdmin.switchOrg(s.otherOrg.ID)
	s.otherAdminInfo, err = s.apiClient.createUser(createUserRequest{
		Email:     "dborry+e2e_other_learner@learntowin.com",
		FirstName: "David",
		LastName:  "Borry",
		Roles:     []string{"ROLE_ADMIN"},
	}, s.superAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(s.otherAdminInfo.ID)

	s.otherAdmin, err = userLogin(s.otherAdminInfo.Email, "password", s.otherOrg.ID, true)
	s.Require().Nil(err)

}

func (s *MainSuite) TearDownSuite() {
	for _, o := range []*organization{s.org, s.otherOrg} {
		_, err := s.db.Exec("delete from organization.organization where slug = ?", o.Slug)
		s.Assert().Nil(err)
		err = s.keycloak.deleteRealm(o.ID)
		s.Assert().Nil(err)
	}

}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(MainSuite))
}

func ptr[T any](t T) *T {
	return &t
}
