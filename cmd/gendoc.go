package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-chi/docgen"
	"github.com/spf13/cobra"

	"github.com/sajib-hassan/warden/internal/api"
)

var (
	routes   bool
	jsonFile bool
)

// gendocCmd represents the gendoc command
var gendocCmd = &cobra.Command{
	Use:   "gendoc",
	Short: "Generate project documentation",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if routes {
			genRoutesDoc()
		}
	},
}

func init() {
	RootCmd.AddCommand(gendocCmd)

	gendocCmd.Flags().BoolVarP(&routes, "routes", "r", false, "create api routes to file")
	gendocCmd.Flags().BoolVarP(&jsonFile, "jsonFile", "j", false,
		"create api routes JSON file otherwise markdown file")
}

func genRoutesDoc() {
	api, _ := api.New()
	if jsonFile {
		fmt.Print("generating routes json file: ")
		jsonapi := docgen.JSONRoutesDoc(api)
		if err := ioutil.WriteFile("routes.json", []byte(jsonapi), 0644); err != nil {
			log.Println(err)

			return
		}
	} else {
		fmt.Print("generating routes markdown file: ")
		md := docgen.MarkdownRoutesDoc(api, docgen.MarkdownOpts{
			ProjectPath: "github.com/sajib-hassan/warden",
			Intro:       "Warden REST API.",
		})
		if err := ioutil.WriteFile("routes.md", []byte(md), 0644); err != nil {
			log.Println(err)
			return
		}
	}
	fmt.Println("OK")
}
