package model

type Invitation struct {
	ID                string `json:"id"`
	InvitedUserID     string `json:"invitedUserId"`
	DownloadedOffline bool   `json:"downloadedOffline"`
}

type InvitationsFilter struct {
	CourseID      string
	InvitedUserID string
}

type InvitationEnrollRequest struct {
	InvitationID string `json:"invitationId"`
}
type InvitationEnrollResponse struct {
	ID           string `json:"id"`
	InvitationID string `json:"invitationId"`
	Status       string `json:"status"`
}

type CardEnrollment struct {
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
type LearningItemEnrollment struct {
	LearningItemEnrollmentId string            `json:"learningItemEnrollmentId,omitempty"`
	CourseEnrollmentId       string            `json:"courseEnrollmentId,omitempty"`
	LearningItemId           string            `json:"learningItemId,omitempty"`
	DeviceId                 string            `json:"deviceId,omitempty"`
	StartedAt                string            `json:"startedAt,omitempty"`
	UpdatedAt                string            `json:"updatedAt,omitempty"`
	Progress                 int32             `json:"progress,omitempty"`
	CardEnrollments          []*CardEnrollment `json:"cardEnrollments,omitempty"`
}
