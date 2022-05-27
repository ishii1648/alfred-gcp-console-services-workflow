package searchers

import (
	"context"
	"fmt"

	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/caching"
	"google.golang.org/api/sqladmin/v1"
)

const cloudSQLEndpoint = "https://console.cloud.google.com/sql/instances/"

type CloudSQLSearcher struct {
	c *sqladmin.Service
}

func (s *CloudSQLSearcher) Search(ctx context.Context, wf *aw.Workflow, fullQuery string, gcpProject string, forceFetch bool) ([]*SearchResult, error) {
	var searchResultList []*SearchResult
	c, err := sqladmin.NewService(ctx)
	if err != nil {
		return searchResultList, err
	}
	s = &CloudSQLSearcher{c: c}

	instances := caching.LoadCloudSQLInstanceListFromCache(wf, ctx, getCurrentFilename(), s.fetch, forceFetch, fullQuery, gcpProject)
	searchResultList = s.getSearchResultList(instances, gcpProject)

	return searchResultList, nil
}

func (s *CloudSQLSearcher) fetch(ctx context.Context, gcpProject string) ([]*sqladmin.DatabaseInstance, error) {
	instances, err := s.c.Instances.List(gcpProject).Do()
	if err != nil {
		return nil, err
	}

	return instances.Items, nil
}

func (s *CloudSQLSearcher) getSearchResultList(instances []*sqladmin.DatabaseInstance, gcpProject string) []*SearchResult {
	var searchResultList []*SearchResult
	for _, instance := range instances {
		searchResult := &SearchResult{
			Title:    instance.Name,
			Subtitle: fmt.Sprintf("%s %s", s.getStatusEmoji(instance.Settings.ActivationPolicy), instance.GceZone),
			Arg:      fmt.Sprintf("%s/%s?overview?project=%s", cloudSQLEndpoint, instance.Name, gcpProject),
		}
		searchResultList = append(searchResultList, searchResult)
	}
	return searchResultList
}

func (s *CloudSQLSearcher) getStatusEmoji(activationPolicy string) string {
	if activationPolicy == "NEVER" {
		return "🔴"
	} else {
		return "🟢"
	}
}
