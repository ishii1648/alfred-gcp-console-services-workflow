package searchers

import (
	"context"
	"fmt"

	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/caching"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/gcp"
	"google.golang.org/api/run/v1"
)

type CloudRunServicesSearcher struct{}

func (s *CloudRunServicesSearcher) Search(ctx context.Context, wf *aw.Workflow, fullQuery string, gcpProject string, gcpService gcp.GcpService, forceFetch bool) error {
	cacheName := getCurrentFilename()
	services := caching.LoadRunServiceListFromCache(wf, ctx, cacheName, s.fetch, forceFetch, fullQuery, gcpProject)

	for _, service := range services {
		if location, ok := service.Metadata.Labels["cloud.googleapis.com/location"]; ok {
			s.addToWorkflow(wf, location, service.Metadata.Name, gcpService, gcpProject)
		}
	}
	return nil
}

func (s *CloudRunServicesSearcher) fetch(ctx context.Context, gcpProject string) ([]run.Service, error) {
	var serviceList []run.Service

	service, err := run.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := service.Namespaces.Services.List(fmt.Sprintf("namespaces/%s", gcpProject)).Do()
	if err != nil {
		return nil, err
	}

	for _, item := range resp.Items {
		serviceList = append(serviceList, *item)
	}

	return serviceList, nil
}

func (s *CloudRunServicesSearcher) addToWorkflow(wf *aw.Workflow, location string, serviceName string, gcpService gcp.GcpService, gcpProject string) {
	wf.NewItem(serviceName).
		Valid(true).
		Var("action", "open-url").
		Arg(fmt.Sprintf("https://console.cloud.google.com/run/detail/%s/%s/metrics", location, serviceName)).
		Icon(&aw.Icon{Value: gcpService.GetIcon()})
}
