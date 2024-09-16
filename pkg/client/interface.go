package client

import "platform_e2e/pkg/client/model"

type APIClient interface {
	// Login - Password login to a specific organization
	Login(username, password, orgID string, acceptTerms bool) (*model.UserCredentials, error)
	// SwitchOrg - Update a user session to switch to a different org
	SwitchOrg(orgID string, credentials *model.UserCredentials) (err error)

	// CreateOrganization - Create a new organization on platform
	CreateOrganization(req *model.CreateOrganizationRequest, credentials *model.UserCredentials) (*model.Organization, error)
	// CreateOrgAttribute - Add a new attribute to an organization
	CreateOrgAttribute(req *model.CreateOrgAttributeRequest, credentials *model.UserCredentials) error
	// OrgAttributes - List all attributes of a given organization
	OrgAttributes(credentials *model.UserCredentials) ([]*model.OrgAttribute, error)

	// CreateUser - Create a new user account in an organization
	CreateUser(req *model.CreateUserRequest, credentials *model.UserCredentials) (*model.User, error)
	// AssignUserAttributes - Assign an attribute value to a given user
	AssignUserAttribute(req *model.AssignUserAttributeRequest, credentials *model.UserCredentials) error

	/////////////////////// Learning ////////////////////
	// CreateLearningPlan - Create a new learning plan
	CreateLearningPlan(req *model.CreateLearningPlanRequest, credentials *model.UserCredentials) (*model.LearningPlan, error)
	// ActivateLearningPlan - Activate a learning plan
	ActivateLearningPlan(learningPlanID string, credentials *model.UserCredentials) (*model.LearningPlan, error)
	// CreateLearningGroup - Create a new learning group
	CreateLearningGroup(req *model.CreateLearningGroupRequest, credentials *model.UserCredentials) (*model.LearningGroup, error)
	// AddGroupsToLearningPlan - Add one or multiple groups to a learning plan
	AddGroupsToLearningPlan(req *model.AddGroupsToLearningPlanRequest, credentials *model.UserCredentials) (*model.LearningPlan, error)
	// CreateCourse - Create a new course
	CreateCourse(req *model.CreateCourseRequest, credentials *model.UserCredentials) (*model.Course, error)
	// ActivateCourse - Activate a course
	ActivateCourse(courseID string, credentials *model.UserCredentials) (*model.Course, error)
	// AddCoursesToLearningPlan - Add one or multiple courses to a learning plan
	AddCoursesToLearningPlan(req *model.AddCoursesToLearningPlanRequest, credentials *model.UserCredentials) (*model.LearningPlan, error)
	// CreateLearningItem - Add a new learning item to a course
	CreateLearningItem(req *model.CreateLearningItemRequest, credentials *model.UserCredentials) (*model.LearningItem, error)
	// CreateCard - Add a new card to a learning item
	CreateCard(req *model.CreateCardRequest, credentials *model.UserCredentials) (*model.Card, error)
	// UpdateCard - Partially updates a card
	UpdateCard(req *model.UpdateCardRequest, credentials *model.UserCredentials) (*model.Card, error)
	// CreateCardsFromFile - Bulk create cards from a json file
	CreateCardsFromFile(learningItemID string, cardsFilepath string, credentials *model.UserCredentials) ([]*model.Card, error)
	// LearningItemCards - Lists all the cards for a given learning item
	LearningItemCards(learningItemID string, credentials *model.UserCredentials) ([]*model.Card, error)
	// CreateMedia - Create a new media item and returns a temporary upload url
	CreateMedia(req *model.CreateMediaRequest, credentials *model.UserCredentials) (*model.Media, error)

	/////////////////////// Enrollments ////////////////////
	// Invitations - Lists user invitations for a given filter
	Invitations(filter *model.InvitationsFilter, credentials *model.UserCredentials) ([]*model.Invitation, error)

	/////////////////////// Offline ////////////////////
	// CourseBundle - Starts a course bundle job for offline
	CourseBundle(req *model.CourseBundleRequest, credentials *model.UserCredentials) (*model.CourseBundleResponse, error)
	// CourseBundleURL - Checks the status of a bundle job and returns download url if completed
	CourseBundleURL(jobID string, credentials *model.UserCredentials) (*model.CourseBundleURLResponse, error)
	// InvitationEnroll - Enroll logged user to a specific invitation
	InvitationEnroll(req *model.InvitationEnrollRequest, credentials *model.UserCredentials) (*model.InvitationEnrollResponse, error)
	// EnrollmentJob - Checks the status of an enrollment job
	EnrollmentJob(jobID string, credentials *model.UserCredentials) (*model.InvitationEnrollResponse, error)
	// CloneEnrollment - For a given invitation id, clone an enrollment for offline
	CloneEnrollment(req *model.CloneEnrollmentRequest, credentials *model.UserCredentials) (*model.CloneEnrollmentResponse, error)
	// SyncEnrollments - Update enrollments on offline database from local device
	SyncEnrollments(req *model.SyncEnrollmentRequest, credentials *model.UserCredentials) ([]*model.SyncEnrollmentsResult, error)
	// DuplicateEnrollments - Get duplicate enrollment from different devices if there are any
	DuplicateEnrollments(req *model.DuplicateEnrollmentsRequest, credentials *model.UserCredentials) ([]*model.LearningItemEnrollment, error)
	// ApproveEnrollment - Approve a learning item enrollment id and resolve potential conflicts between devices
	ApproveEnrollment(req *model.ApproveEnrollmentRequest, credentials *model.UserCredentials) (bool, error)
	// SubmitEnrollments - Submit approved offline enrollments to online platform
	SubmitEnrollments(req *model.SubmitEnrollmentsRequest, credentials *model.UserCredentials) ([]*model.EnrollmentStatus, error)
}
