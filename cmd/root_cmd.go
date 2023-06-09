package cmd

import (
	"micro/config"
	"micro/internal/app"

	"github.com/spf13/cobra"
)

var (
	Runner     CommandLine = &command{}
	configFile             = ""
	debug      bool
)

type CommandLine interface {
	RootCmd() *cobra.Command
	Migrate(cmd *cobra.Command, args []string)
	Seed(cmd *cobra.Command, args []string)
	TestCMD(cmd *cobra.Command, args []string)
}

type command struct {
}

// rootCmd will run the log streamer
var rootCmd = cobra.Command{
	Use:  "micro",
	Long: "A usecase that will validate restful transactions and send them to stripe.",
	Run: func(cmd *cobra.Command, args []string) {
		config.C().Debug = debug
		app.Start()
	},
}

// RootCmd will add flags and subcommands to the different commands
func (c *command) RootCmd() *cobra.Command {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "The configuration file")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "The usecase debug(true is production - false is dev)")

	// add more commands
	rootCmd.AddCommand(&migrateCMD)
	rootCmd.AddCommand(&seedCMD)
	rootCmd.AddCommand(&testCMD)
	return &rootCmd
}
