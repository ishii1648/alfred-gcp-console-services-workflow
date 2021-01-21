package searchers

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/caching"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/gcp"
	"google.golang.org/api/iterator"
)

type PubSubSubscriptionsSearcher struct{}

func (s *PubSubSubscriptionsSearcher) Search(ctx context.Context, wf *aw.Workflow, fullQuery string, gcpProject string, gcpService gcp.GcpService, forceFetch bool) error {
	cacheName := getCurrentFilename()
	subscriptions := caching.LoadGcpPubsubSubscriptionListFromCache(wf, ctx, cacheName, s.fetch, forceFetch, fullQuery, gcpProject)

	for _, sub := range subscriptions {
		s.addToWorkflow(ctx, wf, sub, gcpService, gcpProject)
	}
	return nil
}

func (s *PubSubSubscriptionsSearcher) fetch(ctx context.Context, gcpProject string) ([]*gcp.PubsubSubscription, error) {
	var subscriptions []*gcp.PubsubSubscription
	client, err := pubsub.NewClient(ctx, gcpProject)
	if err != nil {
		return nil, err
	}

	it := client.Subscriptions(ctx)
	for {
		sub, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		subscription := &gcp.PubsubSubscription{
			Name: sub.ID(),
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, nil
}

func (s *PubSubSubscriptionsSearcher) addToWorkflow(ctx context.Context, wf *aw.Workflow, sub *gcp.PubsubSubscription, gcpService gcp.GcpService, gcpProject string) {
	wf.NewItem(sub.Name).
		Valid(true).
		Var("action", "open-url").
		// Subtitle(subtitle).
		Arg(fmt.Sprintf("https://console.cloud.google.com/cloudpubsub/subscription/detail/%s?project=%s", sub.Name, gcpProject)).
		Icon(&aw.Icon{Value: gcpService.GetIcon()})
}
