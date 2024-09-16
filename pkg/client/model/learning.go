package model

type AttributeFilter struct {
	AttributeID    string `json:"attributeId"`
	FilterOperator string `json:"filterOperator"`
	Value          string `json:"value"`
}
type CreateLearningGroupRequest struct {
	Name       string             `json:"name"`
	Attributes []*AttributeFilter `json:"attributes"`
}

type LearningGroup struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	UserCount int    `json:"userCount"`
}

type CreateCourseRequest struct {
	OrganizationId string `json:"organizationId"`
	Title          string `json:"title"`
	VersionName    string `json:"versionName"`
}

type Course struct {
	ID string `json:"id"`
}

type CreateLearningItemRequest struct {
	Course        string `json:"course"`
	Description   string `json:"description"`
	Name          string `json:"name"`
	Points        int    `json:"points"`
	SequenceOrder int    `json:"sequenceOrder"`
	State         string `json:"state"`
	Type          string `json:"type"`
}

type LearningItem struct {
	ID                    string `json:"id"`
	LearningItemVersionID string `json:"learningItemVersionId"`
}

type CardContentBlock struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	JSON any    `json:"json,omitempty"`

	MediaID *string `json:"mediaId,omitempty"`
	Name    *string `json:"name,omitempty"`
}
type CardJSON struct {
	Version       string             `json:"version"`
	Description   string             `json:"description"`
	TemplateType  *string            `json:"templateType"`
	ContentBlocks []CardContentBlock `json:"contentBlocks"`
}
type CreateCardRequest struct {
	LearningItem    string  `json:"learningItem,omitempty"`
	Type            string  `json:"type"`
	SubType         *string `json:"subType,omitempty"`
	Title           string  `json:"title"`
	SequenceOrder   int     `json:"sequenceOrder"`
	ConfidenceCheck bool    `json:"confidenceCheck"`

	JSON any `json:"json"`
}
type CreateCardsRequest struct {
	Cards []*CreateCardRequest `json:"cards"`
}

type Card struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Media []Media  `json:"media"`
	JSON  CardJSON `json:"json"`
}

type CreateMediaRequest struct {
	FileName string `json:"fileName"`
	MimeType string `json:"mimeType"`
}

type Media struct {
	ID              string `json:"id"`
	TemporaryPutURL string `json:"temporaryPutUrl"`
	TemporaryGetURL string `json:"temporaryGetUrl"`
}

type CreateLearningPlanRequest struct {
	Name        string `json:"name"`
	ActivatedAt string `json:"activatedAt"`
}

type LearningPlan struct {
	ID             string          `json:"id"`
	Courses        []Course        `json:"courses"`
	LearningGroups []LearningGroup `json:"learningGroups"`
}

type AddCoursesToLearningPlanRequest struct {
	LearningPlanID string   `json:"-"`
	Courses        []string `json:"courses"`
}

type AddGroupsToLearningPlanRequest struct {
	LearningPlanID   string   `json:"-"`
	LearningGroupIDs []string `json:"learningGroupIds"`
}

type UpdateCardRequest struct {
	ID    string   `json:"-"`
	Title string   `json:"title"`
	JSON  any      `json:"json,omitempty"`
	Media []string `json:"media"`
}
