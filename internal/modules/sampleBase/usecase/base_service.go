package usecase

import (
	"fmt"
	"micro/internal/domain/entity"
	repocontract "micro/internal/modules/sampleBase"
	"regexp"
)

type BaseService struct {
	baseRepository repocontract.IBaseRepository
}

func NewBaseService(repo repocontract.IBaseRepository) BaseService {
	return BaseService{
		baseRepository: repo,
	}
}

func (BaseService) Validate(userID string) bool {
	isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString
	return isAlpha(userID)
}

func (b BaseService) Process(m entity.BaseModel1) (entity.BaseModel2, error) {
	result := entity.BaseModel2{Data: fmt.Sprintf("Hello %s - %d", m.UserID, m.Code)}
	if err := b.baseRepository.StoreBaseModel(m); err != nil {
		return result, err
	}
	if err := b.baseRepository.NotifySomeone(result); err != nil {
		return result, err
	}
	return result, nil
}
