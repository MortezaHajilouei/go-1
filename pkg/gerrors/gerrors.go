package gerrors

import (
	"fmt"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/types/known/durationpb"
)

type gstatusBuilder struct {
	code    codes.Code
	message string
	details []protoiface.MessageV1
}

type Violation struct {
	Field       string
	Description string
}

func NewStatus(code codes.Code) gstatusBuilder {
	return gstatusBuilder{
		code:    code,
		message: "an error occured",
		details: []protoiface.MessageV1{},
	}
}

func (e gstatusBuilder) WithMessage(message string) gstatusBuilder {
	e.message = message
	return e
}

func (e gstatusBuilder) WithMessagef(format string, vals ...interface{}) gstatusBuilder {
	e.message = fmt.Sprintf(format, vals...)
	return e
}

func (e gstatusBuilder) AddBadRequest(violations ...Violation) gstatusBuilder {
	fieldViolations := []*errdetails.BadRequest_FieldViolation{}

	for _, v := range violations {
		fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       v.Field,
			Description: v.Description,
		})
	}

	e.details = append(e.details, &errdetails.BadRequest{
		FieldViolations: fieldViolations,
	})
	return e
}

func (e gstatusBuilder) AddFarsi(message string) gstatusBuilder {
	e.details = append(e.details, &errdetails.LocalizedMessage{
		Locale:  "fa",
		Message: message,
	})
	return e
}

func (e gstatusBuilder) AddError(reason, domain string, meta ...map[string]string) gstatusBuilder {
	_meta := map[string]string{}
	if len(meta) == 1 {
		_meta = meta[0]
	}
	e.details = append(e.details, &errdetails.ErrorInfo{
		Reason:   reason,
		Domain:   domain,
		Metadata: _meta,
	})
	return e
}

func (e gstatusBuilder) AddDebugInfo(detail string, stack []string) gstatusBuilder {
	e.details = append(e.details, &errdetails.DebugInfo{
		Detail:       detail,
		StackEntries: stack,
	})
	return e
}

func (e gstatusBuilder) AddRequestInfo(requestID string) gstatusBuilder {
	e.details = append(e.details, &errdetails.RequestInfo{
		RequestId: requestID,
	})
	return e
}

func (e gstatusBuilder) AddRetryInfo(minDelay time.Duration) gstatusBuilder {
	e.details = append(e.details, &errdetails.RetryInfo{
		RetryDelay: durationpb.New(minDelay),
	})
	return e
}

func (e gstatusBuilder) AddPreconditionFailure(_type, subject, description string) gstatusBuilder {
	e.details = append(e.details, &errdetails.PreconditionFailure{
		Violations: []*errdetails.PreconditionFailure_Violation{{
			Type:        _type,
			Subject:     subject,
			Description: description,
		}},
	})
	return e
}

func (e gstatusBuilder) Make() *status.Status {
	st := status.New(e.code, e.message)
	st, err := st.WithDetails(e.details...)
	if err != nil {
		panic(err)
	}
	return st
}

func (e gstatusBuilder) MakeError() error {
	return e.Make().Err()
}
