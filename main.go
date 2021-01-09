package main

import (
	"flag"
	"log"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/sample-alfred-workflow/workflow"
)

var (
	wf      *aw.Workflow
	query   string
	ymlPath string
)

func init() {
	flag.StringVar(&query, "query", "", "query to use")
	flag.StringVar(&ymlPath, "yml_path", "console-services.yml", "query to use")
	flag.Parse()
	wf = aw.New()
}

func main() {
	wf.Run(func() {
		log.Printf("running workflow with query: `%s`", query)
		query = strings.TrimLeft(query, " ")

		workflow.Run(wf, query, ymlPath)
	})
}
