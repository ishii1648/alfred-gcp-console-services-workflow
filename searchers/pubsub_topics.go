package searchers

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/gcp"
	"google.golang.org/api/iterator"
)

type PubSubTopicsSearcher struct {
	gcpProject string
	gcpService gcp.GcpService
}

func (s *PubSubTopicsSearcher) Search(ctx context.Context, wf *aw.Workflow, query string, gcpProject string, gcpService gcp.GcpService, forceFetch bool) error {
	s = &PubSubTopicsSearcher{
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

func (s *PubSubTopicsSearcher) fetch(ctx context.Context) ([]*pubsub.Topic, error) {
	var topics []*pubsub.Topic
	client, err := pubsub.NewClient(ctx, s.gcpProject)
	if err != nil {
		return nil, err
	}

	it := client.Topics(ctx)
	for {
		t, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}
	return topics, nil
}

func (s *PubSubTopicsSearcher) addToWorkflow(wf *aw.Workflow, topic *pubsub.Topic) {
	wf.NewItem(topic.ID()).
		Valid(true).
		Var("action", "open-url").
		Arg(fmt.Sprintf("https://console.cloud.google.com/cloudpubsub/topic/detail/%s?project=%s", topic.ID(), s.gcpProject)).
		Icon(&aw.Icon{Value: s.gcpService.GetIcon()})
}
