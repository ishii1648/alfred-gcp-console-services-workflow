package searchers

import (
	"context"
	"fmt"

	aw "github.com/deanishe/awgo"
)

const HomeDashboardEndpoint = "https://console.cloud.google.com/home/dashboard"

type ProjectGetSearcher struct{}

func (s *ProjectGetSearcher) Search(ctx context.Context, wf *aw.Workflow, fullQuery string, gcpProject string, forceFetch bool) ([]*SearchResult, error) {
	var searchResultList []*SearchResult

	searchResult := &SearchResult{
		Title: gcpProject,
		Arg:   fmt.Sprintf("%s?project=%s", HomeDashboardEndpoint, gcpProject),
	}

	searchResultList = append(searchResultList, searchResult)

	return searchResultList, nil
}
