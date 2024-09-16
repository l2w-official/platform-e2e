package model

type CourseBundleRequest struct {
	OrgID          string `json:"orgId"`
	CourseID       string `json:"courseId"`
	LearningPlanID string `json:"learningplanId"`
	DeviceID       string `json:"deviceId"`
	OfflineMode    string `json:"offlineMode"`
}

type CourseBundleResponse struct {
	JobID string `json:"jobId"`
}

type CourseBundleURLResponse struct {
	CourseURL    string `json:"courseUrl"`
	BundleStatus string `json:"bundleStatus"`
}

type CloneEnrollmentRequest struct {
	InvitationID string `json:"invitationId"`
}
type CloneEnrollmentResponse struct {
	CourseEnrollmentID      string                    `json:"courseEnrollmentId"`
	CourseID                string                    `json:"courseId"`
	InvitationID            string                    `json:"invitationId"`
	LearningItemEnrollments []*LearningItemEnrollment `json:"learningItemEnrollments"`
	UserID                  string                    `json:"userId"`
}

type SyncEnrollmentRequest struct {
	LearningItemEnrollments []*LearningItemEnrollment `json:"LearningItemEnrollments"`
}

type SyncEnrollmentsResult struct {
	LearningItemEnrollmentID string `json:"learningItemEnrollmentId"`
	Success                  bool   `json:"success"`
	Message                  string `json:"message"`
}

type DuplicateEnrollmentsRequest struct {
	LearningItemEnrollmentIDs []string `json:"learningItemEnrollmentIds"`
}

type ApproveEnrollmentRequest struct {
	LearningItemEnrollmentId string `json:"LearningItemEnrollmentId"`
	DeviceID                 string `json:"deviceId"`
}

type SubmitEnrollmentsRequest struct {
	InvitationIDs []string `json:"invitationIds"`
}
type EnrollmentStatus struct {
	InvitationID string `json:"invitationId"`
	Status       string `json:"status"`
	Message      string `json:"message"`
}
