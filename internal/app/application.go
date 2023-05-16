package app

import (
	"context"
	"fmt"
	"log"
	"micro/adapter/elk"
	"micro/adapter/postgres"
	"micro/internal/middleware"
	"micro/internal/modules/sampleBase/delivery"
	"micro/internal/modules/sampleBase/repository"
	"micro/internal/modules/sampleBase/usecase"
	"micro/internal/server"
	"micro/pkg/logger"
	"os"
	"time"

	"go.uber.org/fx"
)

// StartApplication func
func Start() {
	fmt.Println("\n\n--------------------------------")
	// if go code crashed we get error and line
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// init configs
	for {
		fxNew := fx.New(
			//		fx.Provide(broker.NewNats),
			//		fx.Provide(redis.NewRedis),
			//		fx.Provide(postgres.NewPostgres),
			fx.Provide(postgres.NewGorm),
			fx.Provide(elk.NewLogStash),
			fx.Provide(middleware.NewMiddleware),
			usecase.Module,
			repository.Module,
			delivery.Module,
			fx.Provide(server.New),
			fx.Invoke(logger.InitGlobalLogger),
			fx.Invoke(serve),
		)
		startCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := fxNew.Start(startCtx); err != nil {
			log.Println(err)
			break
		}
		if val := <-fxNew.Done(); val == os.Interrupt {
			break
		}

		stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := fxNew.Stop(stopCtx); err != nil {
			log.Println(err)
			break
		}
	}
}

func serve(lc fx.Lifecycle, server server.IServer) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return server.ListenAndServe()
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown()
		},
	})
}
