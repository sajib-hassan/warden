package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sajib-hassan/warden/internal/db/migrator"
)

var reset bool

// migrateCmd represents the migration command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "use migration tool",
	Long:  `migration uses migrate tool under the hood supporting the same commands and an additional reset command`,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create [-ext E] [-dir D] [-seq] [-digits N] [-format] [-tz] NAME",
	Long: `create [-ext E] [-dir D] [-seq] [-digits N] [-format] [-tz] NAME
	   Create a set of timestamped up/down migrations titled NAME, in directory D with extension E.
	   Use -seq option to generate sequential up/down migrations with N digits.
	   Use -format option to specify a Go time format string. Note: migrations with the same time cause "duplicate migration version" error.
	   Use -tz option to specify the timezone that will be used when generating non-sequential migrations (defaults: UTC).
`,
	Run: func(cmd *cobra.Command, args []string) {
		migrator.ExecuteCreate(args)
	},
}
var gotoCmd = &cobra.Command{
	Use:   "goto",
	Short: "goto V   Migrate to version V`",
	Run: func(cmd *cobra.Command, args []string) {
		migrator.ExecuteGoto(args)
	},
}
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "[N] Apply all or N up migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrator.ExecuteUp(args)
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "[N] [-all]    Apply all or N down migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrator.ExecuteDown(args)
	},
}

var dropCmd = &cobra.Command{
	Use: "drop",
	Short: `drop [-f]    Drop everything inside database
Use -f to bypass confirmation`,
	Run: func(cmd *cobra.Command, args []string) {
		migrator.ExecuteDrop(args)
	},
}

var forceCmd = &cobra.Command{
	Use:   "force",
	Short: `force V      Set version V but don't run migration (ignores dirty state)`,
	Run: func(cmd *cobra.Command, args []string) {
		migrator.ExecuteForce(args)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: `version      Print current migration version`,
	Run: func(cmd *cobra.Command, args []string) {
		migrator.ExecuteVersion(args)
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(createCmd)
	migrateCmd.AddCommand(gotoCmd)
	migrateCmd.AddCommand(upCmd)
	migrateCmd.AddCommand(downCmd)
	migrateCmd.AddCommand(dropCmd)
	migrateCmd.AddCommand(forceCmd)
	migrateCmd.AddCommand(versionCmd)

	migrateCmd.PersistentFlags().StringP("source", "s", "file://internal/db/migrator/migrations", "Location of the migrations (file://path)")
	viper.BindPFlag("source", migrateCmd.PersistentFlags().Lookup("source"))

	downCmd.Flags().Bool("all", false, "[-all] Apply all")
	viper.BindPFlag("applyAll", downCmd.Flags().Lookup("all"))

	dropCmd.Flags().BoolP("force", "f", false, "Use -f to bypass confirmation")
	viper.BindPFlag("forceDrop", dropCmd.Flags().Lookup("force"))

	createCmd.Flags().StringP("ext", "e", "json", "File extension")
	viper.BindPFlag("extPtr", createCmd.Flags().Lookup("ext"))

	createCmd.Flags().StringP("dir", "", "internal/db/migrator/migrations",
		"Directory to place file in",
	)
	viper.BindPFlag("dirPtr", createCmd.Flags().Lookup("dir"))

	createCmd.Flags().StringP("format", "f", "20060102150405",
		`The Go time format string to use. If the string "unix" or "unixNano" is specified, 
				then the seconds or nanoseconds since January 1, 1970 UTC respectively will be used. 
				Caution, due to the behavior of time.Time.Format(), invalid format strings will not error`,
	)
	viper.BindPFlag("formatPtr", createCmd.Flags().Lookup("format"))

	createCmd.Flags().StringP("tz", "t", "UTC",
		`The timezone that will be used for generating timestamps (default: utc)`,
	)
	viper.BindPFlag("timezoneName", createCmd.Flags().Lookup("tz"))

	createCmd.Flags().BoolP("seq", "q", false,
		"Use sequential numbers instead of timestamps (default: false)",
	)
	viper.BindPFlag("seqNumber", createCmd.Flags().Lookup("seq"))

	createCmd.Flags().IntP("digits", "d", 6,
		"The number of digits to use in sequences (default: 6)",
	)
	viper.BindPFlag("seqDigits", createCmd.Flags().Lookup("digits"))

}
