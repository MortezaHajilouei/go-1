package repository

import (
	"context"
	"log"
	"strconv"
	"time"

	"go.uber.org/fx"
	"gorm.io/gorm"

	"micro/adapter/redis"
	"micro/config"
	"micro/internal/domain/entity"
	repocontract "micro/internal/modules/sampleBase"
)

type BaseRepository struct {
	redis redis.Store
	db    *gorm.DB
}

type BaseRepositoryParams struct {
	fx.In

	// Nats  broker.NatsBroker
	// Redis redis.Store
	DB *gorm.DB
}

func NewBaseRepository(params BaseRepositoryParams) repocontract.IBaseRepository {
	return &BaseRepository{
		// nats:  params.Nats,
		// redis: params.Redis,
		// db:    params.DB,
	}
}

func (b *BaseRepository) StoreBaseModel(m entity.BaseModel1) error {
	res := b.db.Create(&m)
	if res.Error != nil {
		log.Println(res.Error)
		return res.Error
	} else {
		log.Println(res)
	}

	return b.redis.Set(context.TODO(),
		config.C().Service.Name+":"+m.UserID,
		strconv.Itoa(int(m.Code)), time.Second*90)
}

func (b *BaseRepository) NotifySomeone(m entity.BaseModel2) error {
	return nil
}
