package searchers

import (
	"context"
	"fmt"

	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/caching"
	"google.golang.org/api/run/v1"
)

const runServiceEndpoint = "https://console.cloud.google.com/run/detail"

type CloudRunServicesSearcher struct {
	c *run.APIService
}

func (s *CloudRunServicesSearcher) Search(ctx context.Context, wf *aw.Workflow, fullQuery string, gcpProject string, forceFetch bool) ([]*SearchResult, error) {
	var searchResultList []*SearchResult
	c, err := run.NewService(ctx)
	if err != nil {
		return nil, err
	}
	s = &CloudRunServicesSearcher{c: c}

	services := caching.LoadRunServiceListFromCache(wf, ctx, getCurrentFilename(), s.fetch, forceFetch, fullQuery, gcpProject)
	searchResultList = s.getSearchResultList(services, gcpProject)

	return searchResultList, nil
}

func (s *CloudRunServicesSearcher) fetch(ctx context.Context, gcpProject string) ([]run.Service, error) {
	var serviceList []run.Service

	resp, err := s.c.Namespaces.Services.List(fmt.Sprintf("namespaces/%s", gcpProject)).Do()
	if err != nil {
		return nil, err
	}

	for _, item := range resp.Items {
		serviceList = append(serviceList, *item)
	}

	return serviceList, nil
}

func (s *CloudRunServicesSearcher) getSearchResultList(services []run.Service, gcpProject string) []*SearchResult {
	var searchResultList []*SearchResult
	for _, service := range services {
		if location, ok := service.Metadata.Labels["cloud.googleapis.com/location"]; ok {
			searchResult := &SearchResult{
				Title: service.Metadata.Name,
				Arg:   fmt.Sprintf("%s/%s/%s?project=%s", runServiceEndpoint, location, service.Metadata.Name, gcpProject),
			}
			searchResultList = append(searchResultList, searchResult)
		}
	}
	return searchResultList
}
