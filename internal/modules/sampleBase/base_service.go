package base

import "micro/internal/domain/entity"

type IBaseService interface {
	Validate(string) bool
	Process(entity.BaseModel1) (entity.BaseModel2, error)
}
