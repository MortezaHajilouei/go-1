package controller

import (
	"micro/api/pb/base"
	"micro/pkg/gerrors"
	"regexp"
)

func ValidateSampleRequest(req *base.SampleRequest) (bool, []gerrors.Violation) {
	violations := gerrors.NewViolationBuilder()
	if req.GetData() > 100 {
		violations.Add("data", "data should be less than or equal to 100")
	}
	if req.GetData() < 0 {
		violations.Add("data", "data should be zero or positive")
	}
	if !regexp.MustCompile(`^[A-Za-z]+$`).MatchString(req.GetUserID()) {
		violations.Add("user_id", "user ID should not contain numbers")
	}
	return violations.IsEmpty(), violations.Make()
}
