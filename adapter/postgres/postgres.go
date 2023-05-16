package postgres

import (
	"context"
	"log"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"

	"micro/config"
)

// NewGorm method job is connect to postgres database and check migration
func NewGorm(lc fx.Lifecycle) *gorm.DB {
	//var err error
	db := new(gorm.DB)
	dsn := "postgresql://" + config.C().Postgres.Username + ":" + config.C().Postgres.Password + "@" + config.C().Postgres.Host + "/" + config.C().Postgres.Schema
	_ = "host=" + config.C().Postgres.Host + " user=" + config.C().Postgres.Username + " password=" + config.C().Postgres.Password + " dbname=" + config.C().Postgres.Schema + " port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			lg := zapgorm2.New(zap.L())
			lg.SetAsDefault()
			_db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
				//Logger: lg,
			})

			if config.C().Debug {
				_db = _db.Debug()
			}

			if err != nil {
				return err
			}

			*db = *_db
			log.Printf("postgres database loaded successfully \n")
			return nil
		},

		OnStop: func(c context.Context) error {
			log.Printf("postgres database connection closed \n")
			DB, _ := db.DB()
			return DB.Close()
		},
	})
	return db
}
