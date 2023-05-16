package cmd

import (
	"context"

	"micro/pkg/it"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var testCMD = cobra.Command{
	Use:     "test",
	Long:    "migrate database strucutures. This will migrate tables",
	Aliases: []string{"m"},
	Run:     Runner.TestCMD,
}

// migrate database with fake data
func (c *command) TestCMD(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	it.MustR(grpc.DialContext(ctx, "localhost:8083", grpc.WithInsecure()))
}
