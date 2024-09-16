package quiz_test

import (
	"database/sql"
	"fmt"
	"platform_e2e/pkg/client"
	"platform_e2e/pkg/client/model"
	"platform_e2e/pkg/config"
	"platform_e2e/pkg/keycloak"
	"testing"

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

	org      *model.Organization
	course   *model.Course
	orgAdmin *model.UserCredentials

	orgAdminInfo *model.User
	learnerInfo  map[string]*model.User

	quiz *model.LearningItem

	learningGroup *model.LearningGroup
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

var firstNames = []string{"John", "Audrey", "David"}
var lastNames = []string{"Doe", "McLean", "Borry"}

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

	s.learnerInfo = make(map[string]*model.User)
	for i := 0; i < 3; i++ {
		email := fmt.Sprintf("dborry+e2e_learner%d@learntowin.com", i)
		s.learnerInfo[email], err = s.apiClient.CreateUser(&model.CreateUserRequest{
			Email:     email,
			FirstName: firstNames[i],
			LastName:  lastNames[i],
			Roles:     []string{"ROLE_LEARNER"},
		}, s.superAdmin)
		s.Require().Nil(err)
		s.Require().NotEmpty(s.learnerInfo[email].ID)
	}
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

	for _, learner := range s.learnerInfo {
		err = s.apiClient.AssignUserAttribute(&model.AssignUserAttributeRequest{
			UserID:      learner.ID,
			AttributeID: attributeID,
			Value:       "Blue",
		}, s.orgAdmin)
		s.Require().Nil(err)
	}

	s.learningGroup, err = s.apiClient.CreateLearningGroup(&model.CreateLearningGroupRequest{
		Name: "Test LG",
		Attributes: []*model.AttributeFilter{
			{AttributeID: attributeID, FilterOperator: "EQ", Value: "Blue"},
		},
	}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().Equal(3, s.learningGroup.UserCount)

}

func (s *MainSuite) setupCourse() {
	var err error
	s.course, err = s.apiClient.CreateCourse(&model.CreateCourseRequest{OrganizationId: s.org.ID, VersionName: "Geography", Title: "Geography"}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().NotEmpty(s.course.ID)

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

}

func (s *MainSuite) TearDownSuite() {
	_, err := s.db.Exec("delete from organization.organization where slug = ?", s.org.Slug)
	s.Assert().Nil(err)
	err = s.keycloak.DeleteRealm(s.org.ID)
	s.Assert().Nil(err)

}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(MainSuite))
}

func ptr[T any](t T) *T {
	return &t
}
