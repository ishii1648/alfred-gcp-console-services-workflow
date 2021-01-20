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

type PubSubTopicsSearcher struct{}

func (s *PubSubTopicsSearcher) Search(ctx context.Context, wf *aw.Workflow, fullQuery string, gcpProject string, gcpService gcp.GcpService, forceFetch bool) error {
	cacheName := getCurrentFilename()
	topics := caching.LoadGcpPubsubTopicListFromCache(wf, ctx, cacheName, s.fetch, forceFetch, fullQuery, gcpProject)

	for _, topic := range topics {
		s.addToWorkflow(wf, topic, gcpService, gcpProject)
	}
	return nil
}

func (s *PubSubTopicsSearcher) fetch(ctx context.Context, gcpProject string) ([]*gcp.PubsubTopic, error) {
	var topics []*gcp.PubsubTopic
	client, err := pubsub.NewClient(ctx, gcpProject)
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
		topic := &gcp.PubsubTopic{Name: t.ID()}
		topics = append(topics, topic)
	}
	return topics, nil
}

func (s *PubSubTopicsSearcher) addToWorkflow(wf *aw.Workflow, topic *gcp.PubsubTopic, gcpService gcp.GcpService, gcpProject string) {
	wf.NewItem(topic.Name).
		Valid(true).
		Var("action", "open-url").
		Arg(fmt.Sprintf("https://console.cloud.google.com/cloudpubsub/topic/detail/%s?project=%s", topic.Name, gcpProject)).
		Icon(&aw.Icon{Value: gcpService.GetIcon()})
}
