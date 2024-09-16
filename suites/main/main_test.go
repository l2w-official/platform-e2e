package main_suite_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"platform_e2e/pkg/client/model"
	"platform_e2e/pkg/transport"
	"strings"
	"time"
)

const (
	deviceAID = "device-a-id"
)

type testOptions struct {
	skip bool
}

func (s *MainSuite) it(statement string, test func(), options ...testOptions) {
	var o testOptions
	if len(options) > 0 {
		o = options[0]
	}

	if o.skip {
		return
	}

	test()
}

func (s *MainSuite) TestBundleCourseInvalidCourseID() {
	_, err := s.apiClient.CourseBundle(&model.CourseBundleRequest{
		OrgID:          s.org.ID,
		CourseID:       "not-found",
		LearningPlanID: s.learningPlan.ID,
		DeviceID:       "test-tablet123",
		OfflineMode:    "SHARED",
	}, s.orgAdmin)
	s.httpCode(err, 404, "could not get course info")
}

func (s *MainSuite) TestBundleCourseLearnerFromOtherOrg() {
	err := s.apiClient.SwitchOrg(s.otherOrg.ID, s.superAdmin)
	s.Require().Nil(err)

	for _, user := range []*model.UserCredentials{s.otherAdmin, s.superAdmin} {
		_, err = s.apiClient.CourseBundle(&model.CourseBundleRequest{
			OrgID:          s.org.ID,
			CourseID:       s.course.ID,
			LearningPlanID: s.learningPlan.ID,
			DeviceID:       "test-tablet123",
			OfflineMode:    "SHARED",
		}, user)
		s.httpCode(err, 404, "could not get course info")
	}
}

func (s *MainSuite) TestBundleCourse() {
	s.it("should bundle course", func() {
		courseBundleResp, err := s.apiClient.CourseBundle(&model.CourseBundleRequest{
			OrgID:          s.org.ID,
			CourseID:       s.course.ID,
			LearningPlanID: s.learningPlan.ID,
			DeviceID:       "test-tablet123",
			OfflineMode:    "SHARED",
		}, s.orgAdmin)
		s.Require().Nil(err)

		var resp *model.CourseBundleURLResponse
		for i := 0; i < 10; i += 1 {
			resp, err = s.apiClient.CourseBundleURL(courseBundleResp.JobID, s.orgAdmin)
			s.Require().Nil(err)

			if strings.EqualFold(resp.BundleStatus, "BUNDLE_UPLOAD_COMPLETED") {
				break
			}

			time.Sleep(time.Second * 5)
			i += 1
		}

		s.Require().Equal("BUNDLE_UPLOAD_COMPLETED", resp.BundleStatus)

	})

	invitations, err := s.apiClient.Invitations(&model.InvitationsFilter{
		CourseID:      s.course.ID,
		InvitedUserID: s.learnerInfo.ID,
	}, s.orgAdmin)
	s.Require().Nil(err)
	s.Require().Equal(1, len(invitations))
	s.Require().Equal(s.learnerInfo.ID, invitations[0].InvitedUserID)

	invitationID := invitations[0].ID

	_, err = s.apiClient.InvitationEnroll(&model.InvitationEnrollRequest{InvitationID: "not-found"}, s.orgAdmin)
	s.httpCode(err, http.StatusNotFound, "could not find invitation")

	invitationEnrollResp, err := s.apiClient.InvitationEnroll(&model.InvitationEnrollRequest{InvitationID: invitationID}, s.orgAdmin)
	s.Require().Nil(err)

	for i := 0; i < 5; i += 1 {
		invitationEnrollResp, err = s.apiClient.EnrollmentJob(invitationEnrollResp.ID, s.orgAdmin)
		s.Require().Nil(err)

		if strings.EqualFold(invitationEnrollResp.Status, "ENROLLMENT_COMPLETED") {
			break
		}

		time.Sleep(time.Second * 5)
		i += 1
	}

	s.Require().Equal("ENROLLMENT_COMPLETED", invitationEnrollResp.Status)

	_, err = s.apiClient.CloneEnrollment(&model.CloneEnrollmentRequest{}, s.orgAdmin)
	s.httpCode(err, http.StatusBadRequest, "invitationId is required")

	_, err = s.apiClient.CloneEnrollment(&model.CloneEnrollmentRequest{InvitationID: "invalid"}, s.orgAdmin)
	s.httpCode(err, http.StatusNotFound, "could not find invitation")

	cloneEnrollmentResp, err := s.apiClient.CloneEnrollment(&model.CloneEnrollmentRequest{InvitationID: invitationID}, s.orgAdmin)
	s.Require().Nil(err)

	b, _ := json.Marshal(cloneEnrollmentResp)
	os.WriteFile("./testdata/clonedEnrollmentsPreSync.json", b, 0644)

	learningItemEnrollments := cloneEnrollmentResp.LearningItemEnrollments
	learningItemEnrollmentIDs := make([]string, len(learningItemEnrollments))
	for j, li := range learningItemEnrollments {
		learningItemEnrollmentIDs[j] = li.LearningItemEnrollmentId
		li.DeviceId = deviceAID
		li.UpdatedAt = time.Now().Format("2006-01-02T15:04:05.999Z07:00")
		li.StartedAt = time.Now().Format("2006-01-02T15:04:05.999Z07:00")

		for i, c := range li.CardEnrollments {
			c.DeviceID = deviceAID
			c.UpdatedAt = time.Now().Format("2006-01-02T15:04:05.999Z07:00")
			c.StartedAt = time.Now().Format("2006-01-02T15:04:05.999Z07:00")
			c.CompletedAt = time.Now().Format("2006-01-02T15:04:05.999Z07:00")

			if !strings.EqualFold(li.LearningItemId, s.quiz.ID) {
				continue
			}

			if i > 2 {
				continue
			}

			c.Progress = 1

			if i == 1 {
				c.Answer = []string{"Europe"}
			}

			if i == 2 {
				c.Answer = []string{"France", "Netherlands"}
				c.Confidence = 1
			}

		}
	}
	syncEnrollmentsResp, err := s.apiClient.SyncEnrollments(&model.SyncEnrollmentRequest{LearningItemEnrollments: learningItemEnrollments}, s.orgAdmin)
	s.Require().Nil(err)

	b, _ = json.Marshal(syncEnrollmentsResp)
	os.WriteFile("./testdata/sync.json", b, 0644)

	for _, r := range syncEnrollmentsResp {
		s.Require().True(r.Success, r.Message)
	}

	// Override with same device
	for _, li := range learningItemEnrollments {
		if !strings.EqualFold(li.LearningItemId, s.quiz.ID) {
			continue
		}

		for i, c := range li.CardEnrollments {
			if i == 0 || i == len(li.CardEnrollments)-1 {
				continue
			}

			if i == 1 {
				c.Answer = nil
			}

			if i == 2 {
				c.Answer = []string{"65c5ec81-b1dd-43c4-aa5f-d569168f8316"}
			}

			if i == 3 {
				c.Answer = nil
			}
		}
	}

	syncEnrollmentsResp, err = s.apiClient.SyncEnrollments(&model.SyncEnrollmentRequest{LearningItemEnrollments: learningItemEnrollments}, s.orgAdmin)
	s.Require().Nil(err)

	b, _ = json.Marshal(syncEnrollmentsResp)
	os.WriteFile("./testdata/sync2.json", b, 0644)

	for _, r := range syncEnrollmentsResp {
		s.Require().True(r.Success, r.Message)
	}

	s.it("shuold approve and submit enrollments", func() {
		for _, learningItemEnrollmentID := range learningItemEnrollmentIDs {
			success, err := s.apiClient.ApproveEnrollment(&model.ApproveEnrollmentRequest{
				LearningItemEnrollmentId: learningItemEnrollmentID,
				DeviceID:                 deviceAID}, s.orgAdmin)
			s.Require().Nil(err)
			s.Require().True(success)
		}

		enrollmentStatuses, err := s.apiClient.SubmitEnrollments(&model.SubmitEnrollmentsRequest{InvitationIDs: []string{invitationID}}, s.orgAdmin)
		s.Require().Nil(err)
		s.Require().Equal(1, len(enrollmentStatuses))

		for _, enrollmentStatus := range enrollmentStatuses {
			s.Require().Contains(enrollmentStatus.Message, "Successfully updated", enrollmentStatus.Message)
		}
	}, testOptions{skip: true})

}

func (s *MainSuite) httpCode(err error, code int, msg ...string) {
	var httpErr *transport.HttpError

	s.Require().True(errors.As(err, &httpErr))

	s.Assert().Equal(code, httpErr.Code)

	if len(msg) > 0 {
		s.Assert().Contains(err.Error(), msg[0])
	}
}
