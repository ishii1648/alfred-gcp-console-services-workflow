package searchers

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/caching"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/gcp"
	"google.golang.org/api/iterator"
)

type GcsBrowserSearcher struct{}

func (s *GcsBrowserSearcher) Search(ctx context.Context, wf *aw.Workflow, fullQuery string, gcpProject string, gcpService gcp.GcpService, forceFetch bool) error {
	cacheName := getCurrentFilename()
	topics := caching.LoadStorageBucketAttrsListFromCache(wf, ctx, cacheName, s.fetch, forceFetch, fullQuery, gcpProject)

	for _, topic := range topics {
		s.addToWorkflow(wf, topic, gcpService, gcpProject)
	}
	return nil
}

func (s *GcsBrowserSearcher) fetch(ctx context.Context, gcpProject string) ([]*storage.BucketAttrs, error) {
	var bucketAttrs []*storage.BucketAttrs
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	it := client.Buckets(ctx, gcpProject)
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

func (s *GcsBrowserSearcher) addToWorkflow(wf *aw.Workflow, b *storage.BucketAttrs, gcpService gcp.GcpService, gcpProject string) {
	wf.NewItem(b.Name).
		Valid(true).
		Var("action", "open-url").
		Subtitle(fmt.Sprintf("%s %s %s", b.LocationType, b.Location, b.StorageClass)).
		Arg(fmt.Sprintf("https://console.cloud.google.com/storage/browser/%s?project=%s", b.Name, gcpProject)).
		Icon(&aw.Icon{Value: gcpService.GetIcon()})
}
