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
	subscriptions := caching.LoadPubsubSubscriptionListFromCache(wf, ctx, cacheName, s.fetch, forceFetch, fullQuery, gcpProject)

	for _, sub := range subscriptions {
		s.addToWorkflow(ctx, wf, sub, gcpService, gcpProject)
	}
	return nil
}

func (s *PubSubSubscriptionsSearcher) fetch(ctx context.Context, gcpProject string) ([]*pubsub.Subscription, error) {
	var subscriptions []*pubsub.Subscription
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
		subscriptions = append(subscriptions, sub)
	}
	return subscriptions, nil
}

func (s *PubSubSubscriptionsSearcher) addToWorkflow(ctx context.Context, wf *aw.Workflow, sub *pubsub.Subscription, gcpService gcp.GcpService, gcpProject string) {
	wf.NewItem(sub.ID()).
		Valid(true).
		Var("action", "open-url").
		// Subtitle(subtitle).
		Arg(fmt.Sprintf("https://console.cloud.google.com/cloudpubsub/subscription/detail/%s?project=%s", sub.ID(), gcpProject)).
		Icon(&aw.Icon{Value: gcpService.GetIcon()})
}
