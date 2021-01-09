package workflow

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	aw "github.com/deanishe/awgo"
)

func Run(wf *aw.Workflow, rawQuery string, ymlPath string) {
	gcpServices := ParseConsoleServicesYml(ymlPath)
	parser := NewParser(strings.NewReader(rawQuery))
	query := parser.Parse()
	defer finalize(wf)

	gcpProject, err := GetCurrentGCPProject()
	if err != nil {
		handleAlertMessage(wf, fmt.Sprintf("failed to get gcp project : %v", err))
		return
	}

	if query.IsEmpty() {
		handleEmptyQuery(wf)
		return
	}

	var gcpService *GcpService
	for i := range gcpServices {
		if gcpServices[i].Id == query.ServiceId {
			gcpService = &gcpServices[i]
			break
		}
	}

	var filterQuery string
	if gcpService == nil {
		filterQuery = query.ServiceId
		SearchServices(wf, gcpServices, gcpProject)
	} else {
		if gcpService.SubServices == nil || len(gcpService.SubServices) <= 0 {
			AddServiceToWorkflow(wf, *gcpService, gcpProject)
			return
		}

		var subService *GcpService
		for i := range gcpService.SubServices {
			if gcpService.SubServices[i].Id == query.SubServiceId {
				subService = &gcpService.SubServices[i]
				break
			}
		}

		if subService == nil {
			filterQuery = query.SubServiceId
			SearchSubServices(wf, *gcpService, gcpProject)
		} else {
			AddSubServiceToWorkflow(wf, *gcpService, *subService, gcpProject)
		}
	}

	if filterQuery != "" {
		log.Printf("filtering with query %s", filterQuery)
		res := wf.Filter(filterQuery)
		log.Printf("%d results match %q", len(res), filterQuery)
	}
}

func finalize(wf *aw.Workflow) {
	wf.SendFeedback()
}

func handleEmptyQuery(wf *aw.Workflow) {
	log.Println("no search type parsed")
	wf.NewItem("Search for an GCP Service ...").
		Subtitle("e.g., gke, gcs, cloud run ...")

	if wf.UpdateCheckDue() {
		if err := wf.CheckForUpdate(); err != nil {
			wf.FatalError(err)
		}
	}
}

func SearchServices(wf *aw.Workflow, gcpServices []GcpService, gcpProject string) {
	for i := range gcpServices {
		AddServiceToWorkflow(wf, gcpServices[i], gcpProject)
	}
}

func SearchSubServices(wf *aw.Workflow, gcpService GcpService, gcpProject string) {
	for _, subService := range gcpService.SubServices {
		AddSubServiceToWorkflow(wf, gcpService, subService, gcpProject)
	}
}

func AddServiceToWorkflow(wf *aw.Workflow, gcpService GcpService, gcpProject string) {
	title := gcpService.Id

	subtitle := ""
	if len(gcpService.SubServices) > 0 {
		subtitle += "🗂 "
	}

	subtitle += gcpService.Name
	if gcpService.ShortName != "" {
		subtitle += " (" + gcpService.ShortName + ")"
	}
	subtitle += " – " + gcpService.Description

	wf.NewItem(title).
		Valid(true).
		Var("action", "open-url").
		Subtitle(subtitle).
		Autocomplete(gcpService.Id + " ").
		UID(gcpService.Id).
		Arg(fmt.Sprintf("%s?project=%s", gcpService.Url, gcpProject)).
		Icon(&aw.Icon{Value: GetIcon(&gcpService)})
}

func AddSubServiceToWorkflow(wf *aw.Workflow, gcpService, subService GcpService, gcpProject string) {
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

	wf.NewItem(title).
		Valid(true).
		Var("action", "open-url").
		Subtitle(subtitle).
		Autocomplete(gcpService.Id + " " + subService.Id + " ").
		UID(gcpService.Id).
		Arg(fmt.Sprintf("%s?project=%s", gcpService.Url, gcpProject)).
		Icon(&aw.Icon{Value: GetIcon(&gcpService)})
}

func GetIcon(gcpService *GcpService) string {
	iconPath := fmt.Sprintf("images/%s.png", gcpService.Id)
	if _, err := os.Stat(iconPath); err != nil {
		return "images/gcp.png"
	}
	return iconPath
}

func GetCurrentGCPProject() (string, error) {
	var project string

	gcp_config := os.Getenv("ALFRED_GCP_CONSOLE_SERVICES_WORKFLOW_GCP_CONFIG")
	if gcp_config == "" {
		return project, fmt.Errorf("You should set environment : ALFRED_GCP_CONSOLE_SERVICES_WORKFLOW_GCP_CONFIG")
	}

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

func handleAlertMessage(wf *aw.Workflow, title string) {
	wf.NewItem(title).
		Valid(true).
		Var("action", "open-url").
		Arg("https://github.com/rkoval/alfred-aws-console-services-workflow/blob/master/CONTRIBUTING.md").
		Icon(aw.IconNote)
}
