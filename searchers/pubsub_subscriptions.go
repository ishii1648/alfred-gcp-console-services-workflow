package searchers

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/gcp"
	"google.golang.org/api/iterator"
)

type PubSubSubscriptionsSearcher struct {
	gcpProject string
	gcpService gcp.GcpService
}

func (s *PubSubSubscriptionsSearcher) Search(ctx context.Context, wf *aw.Workflow, gcpProject string, gcpService gcp.GcpService) error {
	s = &PubSubSubscriptionsSearcher{
		gcpProject: gcpProject,
		gcpService: gcpService,
	}

	subscriptions, err := s.fetch(ctx)
	if err != nil {
		return err
	}

	for _, sub := range subscriptions {
		s.addToWorkflow(ctx, wf, sub)
	}
	return nil
}

func (s *PubSubSubscriptionsSearcher) fetch(ctx context.Context) ([]*pubsub.Subscription, error) {
	var subscriptions []*pubsub.Subscription
	client, err := pubsub.NewClient(ctx, s.gcpProject)
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

func (s *PubSubSubscriptionsSearcher) addToWorkflow(ctx context.Context, wf *aw.Workflow, sub *pubsub.Subscription) {
	// conf, err := sub.Config(ctx)
	// if err != nil {
	// 	wf.NewItem("Failed to get subscription config").
	// 		Valid(true).
	// 		Var("action", "open-url").
	// 		Subtitle(err.Error()).
	// 		Icon(aw.IconError)
	// 	return
	// }

	// var subtitle string
	// subtitle = fmt.Sprintf("%s %d", conf.Topic.ID(), conf.RetentionDuration)

	wf.NewItem(sub.ID()).
		Valid(true).
		Var("action", "open-url").
		// Subtitle(subtitle).
		Arg(fmt.Sprintf("https://console.cloud.google.com/cloudpubsub/subscription/detail/%s?project=%s", sub.ID(), s.gcpProject)).
		Icon(&aw.Icon{Value: s.gcpService.GetIcon()})
}
