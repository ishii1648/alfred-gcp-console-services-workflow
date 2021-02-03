package searchers

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/caching"
	"google.golang.org/api/iterator"
)

const gcsBucketEndpoint = "https://console.cloud.google.com/storage/browser"

type GcsBrowserSearcher struct {
	c *storage.Client
}

func (s *GcsBrowserSearcher) Search(ctx context.Context, wf *aw.Workflow, fullQuery string, gcpProject string, forceFetch bool) ([]*SearchResult, error) {
	var searchResultList []*SearchResult
	c, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	s = &GcsBrowserSearcher{c: c}

	bucketAttrs := caching.LoadStorageBucketAttrsListFromCache(wf, ctx, getCurrentFilename(), s.fetch, forceFetch, fullQuery, gcpProject)
	searchResultList = s.getSearchResultList(bucketAttrs, gcpProject)

	return searchResultList, nil
}

func (s *GcsBrowserSearcher) fetch(ctx context.Context, gcpProject string) ([]*storage.BucketAttrs, error) {
	var bucketAttrs []*storage.BucketAttrs

	it := s.c.Buckets(ctx, gcpProject)
	for {
		bucketAttr, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		bucketAttrs = append(bucketAttrs, bucketAttr)
	}
	return bucketAttrs, nil
}

func (s *GcsBrowserSearcher) getSearchResultList(bucketAttrs []*storage.BucketAttrs, gcpProject string) []*SearchResult {
	var searchResultList []*SearchResult
	for _, bucketAttr := range bucketAttrs {
		searchResult := &SearchResult{
			Title:    bucketAttr.Name,
			Subtitle: fmt.Sprintf("%s %s %s", bucketAttr.LocationType, bucketAttr.Location, bucketAttr.StorageClass),
			Arg:      fmt.Sprintf("%s/%s?project=%s", gcsBucketEndpoint, bucketAttr.Name, gcpProject),
		}
		searchResultList = append(searchResultList, searchResult)
	}
	return searchResultList
}
