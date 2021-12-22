package searchers

import (
	"context"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/caching"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"
)

type ProjectSearcher struct {
	svc *cloudresourcemanager.Service
}

func (s *ProjectSearcher) Search(ctx context.Context, wf *aw.Workflow, fullQuery string, gcpProject string, forceFetch bool) ([]*SearchResult, error) {
	var searchResultList []*SearchResult

	c, err := google.DefaultClient(ctx, cloudresourcemanager.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	svc, err := cloudresourcemanager.New(c)
	if err != nil {
		return nil, err
	}
	s = &ProjectSearcher{svc: svc}

	projects := caching.LoadCloudresourcemanagerProjectListFromCache(wf, ctx, getCurrentFilename(), s.fetch, forceFetch, fullQuery, gcpProject)
	searchResultList = s.getSearchResultList(projects)

	return searchResultList, nil
}

func (s *ProjectSearcher) fetch(ctx context.Context, gcpProject string) ([]*cloudresourcemanager.Project, error) {
	resp, err := s.svc.Projects.List().Do()
	if err != nil {
		return nil, err
	}

	return resp.Projects, nil
}

func (s *ProjectSearcher) getSearchResultList(projects []*cloudresourcemanager.Project) []*SearchResult {
	var searchResultList []*SearchResult

	for _, project := range projects {
		// exluce automate created projects
		if strings.HasPrefix(project.ProjectId, "sys-") {
			continue
		}
		if project.LifecycleState != "ACTIVE" {
			continue
		}
		searchResult := &SearchResult{
			Title: project.ProjectId,
		}
		searchResultList = append(searchResultList, searchResult)
	}

	return searchResultList
}
