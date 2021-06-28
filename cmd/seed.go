package cmd

import (
	"github.com/spf13/cobra"

	"github.com/sajib-hassan/warden/internal/db/seeder"
)

// seedCmd to seed database
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "use seeder for seeding database",
	Long:  `Seed database`,
	Run: func(cmd *cobra.Command, args []string) {
		seeder.Execute(args)
	},
}

func init() {
	RootCmd.AddCommand(seedCmd)
}
