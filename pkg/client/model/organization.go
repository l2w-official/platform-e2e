package model

type Organization struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	IDPType     string `json:"idpType"`
	IDPGroupID  string `json:"idpGroupId"`
	IDPClientID string `json:"idpClientId"`
	Status      string `json:"statis"`
}

type AttributeOption struct {
	Label         string `json:"label"`
	SequenceOrder int    `json:"sequenceOrder"`
}
type OrgAttribute struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Type             string             `json:"types"`
	Status           string             `json:"status"`
	AttributeOptions []*AttributeOption `json:"attributeOptions"`
}

type CreateOrganizationRequest struct {
	Slug   string `json:"slug"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type CreateOrganizationResponse struct {
	Organization Organization `json:"organization"`
}

type CreateOrgAttributeRequest struct {
	AttributeOptions []*AttributeOption `json:"attributeOptions"`
	Name             string             `json:"name"`
	Organization     string             `json:"organization"`
	Type             string             `json:"type"`
}

type User struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Roles     []string `json:"roles"`
}

type CreateUserRequest struct {
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Roles     []string `json:"roles"`
}

type AssignUserAttributeRequest struct {
	UserID      string `json:"userId"`
	AttributeID string `json:"attributeId"`
	Value       string `json:"value"`
}
