package cmd

import (
	"fmt"
	"log"
	net "net"

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

	// Here you will define your flags and configuration settings.
	viper.SetDefault("port", "8000")
	viper.SetDefault("log_level", "debug")

	viper.SetDefault("auth_login_url", "http://localhost:3000/login")
	viper.SetDefault("auth_login_pin_length", 5)
	viper.SetDefault("auth_login_token_expiry", "11m")
	viper.SetDefault("auth_jwt_secret", "random")
	viper.SetDefault("auth_jwt_expiry", "15m")
	viper.SetDefault("auth_jwt_refresh_expiry", "1h")

	// Here you will define your flags and configuration settings.
	serveCmd.PersistentFlags().IntP("PORT", "p",
		8000, `Port on which the server will listen. Default port is 8000`,
	)
	viper.BindPFlag("PORT", serveCmd.PersistentFlags().Lookup("PORT"))
}

func runServe(cmd *cobra.Command, args []string) {
	server, err := api.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	server.Start()

}
