package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/docgen"
	"github.com/spf13/cobra"

	"github.com/sajib-hassan/warden/internal/api"
)

var (
	routes   bool
	exportTo string
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
			router, _ := api.New()
			switch exportTo {
			case "json":
				genRoutesJSONDoc(router)
			case "md":
				genRoutesMarkdownDoc(router)
			default:
				printRoutes(router)
			}
		} else {
			fmt.Println("Please call with -r flag")
		}
	},
}

func init() {
	RootCmd.AddCommand(gendocCmd)

	gendocCmd.Flags().BoolVarP(&routes, "routes", "r", false, "create api routes to file")
	gendocCmd.Flags().StringVarP(&exportTo, "export", "e", "cli",
		`create api routes to file as JSON or Markdown format. 
				Options are json | md | cli`,
	)
}

func printRoutes(router *chi.Mux) {
	fmt.Println("Printing available routes: ")
	docgen.PrintRoutes(router)
	fmt.Println("OK")
}

func genRoutesJSONDoc(router *chi.Mux) {
	fmt.Println("generating routes json file: ")
	jsonapi := docgen.JSONRoutesDoc(router)
	if err := ioutil.WriteFile("routes.json", []byte(jsonapi), 0644); err != nil {
		log.Println(err)

		return
	}
	fmt.Println("OK")
}

func genRoutesMarkdownDoc(router *chi.Mux) {
	fmt.Println("generating routes markdown file: ")
	md := docgen.MarkdownRoutesDoc(router, docgen.MarkdownOpts{
		ProjectPath: "github.com/sajib-hassan/warden",
		Intro:       "Warden REST API.",
	})
	if err := ioutil.WriteFile("routes.md", []byte(md), 0644); err != nil {
		log.Println(err)
		return
	}

	fmt.Println("OK")
}
