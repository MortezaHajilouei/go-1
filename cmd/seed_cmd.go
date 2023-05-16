package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
	"log"
)

var (
	seedCMD = cobra.Command{
		Use:  "seed database",
		Long: "seed database strucutures. This will seed tables",
		Run:  Runner.Seed,
	}
)

// seed database with fake data
func (c *command) Seed(cmd *cobra.Command, args []string) {
	fx.New(
		//fx.Provide(elk.NewLogStash),
		//fx.Invoke(logger.InitGlobalLogger),
		fx.Invoke(seed),
	).Start(context.TODO())
}

func seed(lc fx.Lifecycle, db *bun.DB) {
	// Do all your seeding here
	lc.Append(fx.Hook{OnStart: func(c context.Context) error {
		log.Println("Data seeded successfully")
		return nil
	}})
}
