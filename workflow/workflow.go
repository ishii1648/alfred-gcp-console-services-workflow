package workflow

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/gcp"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/searchers"
)

const ProjectCacheFile = "/tmp/alfred-gcp-console-services-workflow/current-project"

type Workflow struct {
	wf *aw.Workflow
}

func Run(wf *aw.Workflow, rawQuery string, ymlPath string, forceFetch bool) {
	log.Println("using workflow cacheDir: " + wf.CacheDir())
	log.Println("using workflow dataDir: " + wf.DataDir())

	workflow := &Workflow{
		wf: wf,
	}

	gcpServices := gcp.ParseConsoleServicesYml(ymlPath)
	parser := NewParser(strings.NewReader(rawQuery))
	query := parser.Parse()
	log.Printf("query parsed : %#v", query)

	defer workflow.finalize()

	if query.IsEmpty() {
		workflow.handleEmptyQuery()
		return
	}

	var gcpProject string

	fileInfo, err := os.Stat(ProjectCacheFile)
	if err == nil && time.Now().Before(fileInfo.ModTime().Add(1*time.Hour)) {
		bGCPProject, err := ioutil.ReadFile(ProjectCacheFile)
		if err != nil {
			workflow.handleAlertMessage(fmt.Sprintf("failed to read ProjectCacheFile : %v", err))
		}
		gcpProject = string(bGCPProject)
	}

	if gcpProject == "" {
		configPath := workflow.getGcpConfigPath()
		gcpProject, err = GetCurrentGCPProject(configPath)
		if err != nil {
			workflow.handleAlertMessage(fmt.Sprintf("failed to get gcp project : %v", err))
			return
		}
	}

	ctx := context.Background()

	var gcpService *gcp.GcpService
	for i := range gcpServices {
		if gcpServices[i].Id == query.ServiceId {
			gcpService = &gcpServices[i]
			break
		}
	}

	var filterQuery string
	if gcpService == nil {
		filterQuery = query.ServiceId
		workflow.SearchServices(gcpServices, gcpProject)
	} else {
		var subService *gcp.GcpService
		for i := range gcpService.SubServices {
			if gcpService.SubServices[i].Id == query.SubServiceId {
				subService = &gcpService.SubServices[i]
				break
			}
		}

		serviceId := query.ServiceId
		if serviceId != "project" && query.SubServiceId != "" {
			serviceId += "_" + query.SubServiceId
		}
		searcher := searchers.SearchersByServiceId[serviceId]
		if searcher != nil {
			filterQuery = query.Filter
			results, err := searcher.Search(ctx, wf, rawQuery, gcpProject, forceFetch)
			if err != nil {
				wf.FatalError(err)
			}
			for _, result := range results {
				workflow.AddSearchedServiceToWorkflow(*gcpService, result.Title, result.Subtitle, result.Arg)
			}
			if serviceId == "project" {
				filterQuery = query.SubServiceId
			}
		} else {
			if subService == nil {
				filterQuery = query.SubServiceId
				workflow.SearchSubServices(*gcpService, gcpProject)
			} else {
				workflow.AddSubServiceToWorkflow(*gcpService, *subService, gcpProject)
			}
		}
	}

	if filterQuery != "" {
		log.Printf("filtering with query %s", filterQuery)
		res := wf.Filter(filterQuery)
		log.Printf("%d results match %q", len(res), filterQuery)
	}
}

func (w *Workflow) finalize() {
	if w.wf.IsEmpty() {
		w.wf.NewItem("No matching services found").
			Subtitle("Try another query (example: `gcp gke clusters`)").
			Icon(aw.IconNote)
	}
	w.wf.SendFeedback()
}

func (w *Workflow) handleEmptyQuery() {
	w.wf.NewItem("Search for an GCP Service ...").
		Subtitle("e.g., gke, gcs, cloud run ...")

	if w.wf.UpdateCheckDue() {
		if err := w.wf.CheckForUpdate(); err != nil {
			w.wf.FatalError(err)
		}
	}
}

func (w *Workflow) SearchServices(gcpServices []gcp.GcpService, gcpProject string) {
	for i := range gcpServices {
		w.AddServiceToWorkflow(gcpServices[i], gcpProject)
	}
}

func (w *Workflow) SearchSubServices(gcpService gcp.GcpService, gcpProject string) {
	if len(gcpService.SubServices) > 0 {
		for _, subService := range gcpService.SubServices {
			w.AddSubServiceToWorkflow(gcpService, subService, gcpProject)
		}
		return
	}
	w.AddServiceToWorkflow(gcpService, gcpProject)
}

func (w *Workflow) AddServiceToWorkflow(gcpService gcp.GcpService, gcpProject string) {
	title := gcpService.Id

	subtitle := ""
	if len(gcpService.SubServices) > 0 {
		subtitle += "🗂 "
	}

	searcher := searchers.SearchersByServiceId[gcpService.Id]
	if searcher != nil {
		subtitle += "🔎 "
	}

	subtitle += gcpService.Name
	if gcpService.ShortName != "" {
		subtitle += " (" + gcpService.ShortName + ")"
	}
	subtitle += " – " + gcpService.Description

	w.wf.NewItem(title).
		Valid(true).
		Var("action", "open-url").
		Subtitle(subtitle).
		Autocomplete(gcpService.Id + " ").
		UID(gcpService.Id).
		Arg(fmt.Sprintf("%s?project=%s", gcpService.Url, gcpProject)).
		Icon(&aw.Icon{Value: gcpService.GetIcon()})
}

func (w *Workflow) AddSubServiceToWorkflow(gcpService, subService gcp.GcpService, gcpProject string) {
	title := gcpService.Id + " " + subService.Id
	subtitle := ""

	if gcpService.ShortName != "" {
		subtitle += gcpService.ShortName + " – "
	} else {
		subtitle += gcpService.GetName() + " – "
	}

	subtitle += subService.Name
	if subService.Description != "" {
		subtitle += " – " + subService.Description
	}

	searcher := searchers.SearchersByServiceId[gcpService.Id+"_"+subService.Id]
	if searcher != nil {
		subtitle += "🔎 "
	}

	w.wf.NewItem(title).
		Valid(true).
		Var("action", "open-url").
		Subtitle(subtitle).
		Autocomplete(gcpService.Id + " " + subService.Id + " ").
		UID(gcpService.Id).
		Arg(fmt.Sprintf("%s?project=%s", subService.Url, gcpProject)).
		Icon(&aw.Icon{Value: gcpService.GetIcon()})
}

func (w *Workflow) AddSearchedServiceToWorkflow(gcpService gcp.GcpService, title, subtitle, arg string) {
	if gcpService.Id == "project" && arg == "" {
		w.wf.NewItem(title).
			Valid(true).
			Var("output", "write-file").
			Autocomplete(title + " ").
			UID(title).
			Arg(title)
	} else {
		w.wf.NewItem(title).
			Valid(true).
			Var("action", "open-url").
			Subtitle(subtitle).
			Arg(arg).
			Icon(&aw.Icon{Value: gcpService.GetIcon()})
	}
}

func (w *Workflow) getGcpConfigPath() string {
	gcp_config := os.Getenv("ALFRED_GCP_CONSOLE_SERVICES_WORKFLOW_GCP_CONFIG")
	if gcp_config == "" {
		cacheDirList := strings.Split(w.wf.CacheDir(), "/")
		gcp_config = fmt.Sprintf("/%s/%s/.config/gcloud/configurations/config_default", cacheDirList[1], cacheDirList[2])
	}
	return gcp_config
}

func (w *Workflow) handleAlertMessage(title string) {
	w.wf.NewItem(title).
		Valid(true).
		Var("action", "open-url").
		Arg("https://github.com/rkoval/alfred-aws-console-services-workflow/blob/master/CONTRIBUTING.md").
		Icon(aw.IconNote)
}

func GetCurrentGCPProject(gcp_config string) (string, error) {
	var project string

	log.Printf("gcp_config : %s", gcp_config)

	f, err := os.Open(gcp_config)
	if err != nil {
		return project, err
	}
	defer f.Close()

	var projectLine string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "project") {
			// project = xxx
			projectLine = scanner.Text()
			break
		}
	}

	if err = scanner.Err(); err != nil {
		return project, err
	}

	if len(projectLine) <= 0 {
		return project, fmt.Errorf("no project")
	}

	project = strings.Split(projectLine, "=")[1]
	project = strings.TrimSpace(project)

	return project, nil
}
