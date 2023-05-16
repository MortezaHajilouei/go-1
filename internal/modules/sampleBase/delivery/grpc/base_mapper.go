package controller

import (
	"micro/api/pb/base"
	"micro/internal/domain/entity"
)

func SampleRequestToBaseModel1(r *base.SampleRequest) entity.BaseModel1 {
	return entity.BaseModel1{
		UserID: r.GetUserID(),
		Code:   r.GetData(),
	}
}

func BaseModel2ToSampleResponse(m entity.BaseModel2) *base.SampleResponse {
	return &base.SampleResponse{
		Data: m.Data,
	}
}
