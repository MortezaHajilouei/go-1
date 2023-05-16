package base

import (
	"micro/internal/domain/entity"
)

type IBaseRepository interface {
	StoreBaseModel(entity.BaseModel1) error
	NotifySomeone(entity.BaseModel2) error
}
