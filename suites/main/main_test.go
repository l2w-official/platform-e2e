package main_suite_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	loadConfig()
}

func TestE2e(t *testing.T) {

	// INIT
	require.NotNil(t, config)

	token, err := login(config.l2wOrgID, "sbengtson@learntowin.com", "password")
	require.Nil(t, err)
	require.NotEmpty(t, token)

	contextToken, err := userContext(config.l2wOrgID, token, false)
	require.Nil(t, err)
	require.NotEmpty(t, contextToken)

	org, err := createOrganization(createOrganizationRequest{Slug: "atg", Name: "atg", Status: "ACTIVE"}, token, contextToken)
	require.Nil(t, err)
	require.Equal(t, org.Slug, "atg")

	contextToken, err = userContext(org.ID, token, false)
	require.Nil(t, err)
	require.NotEmpty(t, contextToken)

	admin, err := createUser(createUserRequest{
		Email:     "dborry+e2e_admin@learntowin.com",
		FirstName: "David",
		LastName:  "Borry",
		Roles:     []string{"ROLE_ADMIN"},
	}, token, contextToken)
	require.Nil(t, err)
	require.NotEmpty(t, admin.ID)

	learner, err := createUser(createUserRequest{
		Email:     "dborry+e2e_learner@learntowin.com",
		FirstName: "David",
		LastName:  "Borry",
		Roles:     []string{"ROLE_LEARNER"},
	}, token, contextToken)
	require.Nil(t, err)
	require.NotEmpty(t, learner.ID)

	// Setup user group
	adminToken, err := login(org.ID, "dborry+e2e_admin@learntowin.com", "password")
	require.Nil(t, err)
	require.NotEmpty(t, adminToken)

	time.Sleep(time.Second * 5)

	adminContextToken, err := userContext(org.ID, adminToken, true)
	require.Nil(t, err)
	require.NotEmpty(t, adminContextToken)

	err = createOrgAttribute(createOrgAttributeRequest{
		AttributeOptions: []*attributeOption{
			{Label: "Red", SequenceOrder: 0},
			{Label: "Green", SequenceOrder: 1},
			{Label: "Blue", SequenceOrder: 2},
		},
		Name:         "Color",
		Organization: fmt.Sprintf("/api/organizations/%s", org.ID),
		Type:         "SINGLE_SELECT",
	}, adminToken, adminContextToken)
	require.Nil(t, err)

	attributes, err := orgAttributes(adminToken, adminContextToken)
	require.Nil(t, err)
	require.Equal(t, 1, len(attributes))

	attributeID := attributes[0].ID
	require.NotEmpty(t, attributeID)

	err = assignUserAttributes(assignUserAttributesRequest{
		UserID:      learner.ID,
		AttributeID: attributeID,
		Value:       "Blue",
	}, adminToken, adminContextToken)
	require.Nil(t, err)

	learningGroup, err := createLearningGroup(createLearningGroupRequest{
		Name: "Test LG",
		Attributes: []*attributeFilter{
			{AttributeID: attributeID, FilterOperator: "EQ", Value: "Blue"},
		},
	}, adminToken, adminContextToken)
	require.Nil(t, err)
	require.Equal(t, 1, learningGroup.UserCount)

}
