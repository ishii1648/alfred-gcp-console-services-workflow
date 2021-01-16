package main

import (
	"flag"
	"log"
	"os"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/workflow"
)

var (
	wf         *aw.Workflow
	forceFetch bool
	query      string
	ymlPath    string
)

func init() {
	flag.BoolVar(&forceFetch, "fetch", false, "force fetch via GCP instead of cache")
	flag.StringVar(&query, "query", "", "query to use")
	flag.StringVar(&ymlPath, "yml_path", "console-services.yml", "query to use")
	flag.Parse()
	wf = aw.New()
}

func main() {
	wf.Run(func() {
		log.Printf("%v", os.Args)
		log.Printf("running workflow with query: `%s`", query)
		query = strings.TrimLeft(query, " ")

		workflow.Run(wf, query, ymlPath, forceFetch)
	})
}
