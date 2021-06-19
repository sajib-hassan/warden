package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sajib-hassan/warden/internal/api"
)

// ServeCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"s"},
	Short:   "Start authorizer API Server",
	Long: `Start authorizer API Server 
with the provided configurations.`,
	Example: `$ go run main.go serve
or
$ go run main.go s`,
	Run: runServe,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		portStr := viper.GetString("PORT")
		listener, err := net.Listen("tcp", ":"+portStr)
		if err != nil {
			return fmt.Errorf("port %s is not available", portStr)
		}

		listener.Close()
		return nil
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().IntP("port", "p",
		8000, `Port on which the server will listen. Default port is 8000`,
	)
	viper.BindPFlag("port", serveCmd.PersistentFlags().Lookup("port"))

	serveCmd.PersistentFlags().StringP("environment", "e",
		"dev", `Running environment dev, stg, prod, test. (default is dev)`,
	)
	viper.BindPFlag("environment", serveCmd.PersistentFlags().Lookup("environment"))
}

func runServe(cmd *cobra.Command, args []string) {
	server, err := api.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	server.Start()

}
