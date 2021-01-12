package searchers

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/gcp"
	"google.golang.org/api/iterator"
)

type GcsSearcher struct {
	gcpProject string
	gcpService gcp.GcpService
}

func (s *GcsSearcher) Search(ctx context.Context, wf *aw.Workflow, gcpProject string, gcpService gcp.GcpService) error {
	s = &GcsSearcher{
		gcpProject: gcpProject,
		gcpService: gcpService,
	}

	topics, err := s.fetch(ctx)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		s.addToWorkflow(wf, topic)
	}
	return nil
}

func (s *GcsSearcher) fetch(ctx context.Context) ([]*storage.BucketAttrs, error) {
	var bucketAttrs []*storage.BucketAttrs
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	it := client.Buckets(ctx, s.gcpProject)
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

func (s *GcsSearcher) addToWorkflow(wf *aw.Workflow, b *storage.BucketAttrs) {
	wf.NewItem(b.Name).
		Valid(true).
		Var("action", "open-url").
		Subtitle(fmt.Sprintf("%s %s %s", b.LocationType, b.Location, b.StorageClass)).
		Arg(fmt.Sprintf("https://console.cloud.google.com/storage/browser/%s?project=%s", b.Name, s.gcpProject)).
		Icon(&aw.Icon{Value: s.gcpService.GetIcon()})
}
