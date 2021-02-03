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

const pubsubTopicEndpoint = "https://console.cloud.google.com/cloudpubsub/topic/detail"

type PubSubTopicsSearcher struct {
	c *pubsub.Client
}

func (s *PubSubTopicsSearcher) Search(ctx context.Context, wf *aw.Workflow, fullQuery string, gcpProject string, forceFetch bool) ([]*SearchResult, error) {
	var searchResultList []*SearchResult
	c, err := pubsub.NewClient(ctx, gcpProject)
	if err != nil {
		return nil, err
	}
	s = &PubSubTopicsSearcher{c: c}

	topics := caching.LoadGcpPubsubTopicListFromCache(wf, ctx, getCurrentFilename(), s.fetch, forceFetch, fullQuery, gcpProject)
	searchResultList = s.getSearchResultList(topics, gcpProject)

	return searchResultList, nil
}

func (s *PubSubTopicsSearcher) fetch(ctx context.Context, gcpProject string) ([]*gcp.PubsubTopic, error) {
	var topics []*gcp.PubsubTopic

	it := s.c.Topics(ctx)
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

func (s *PubSubTopicsSearcher) getSearchResultList(topics []*gcp.PubsubTopic, gcpProject string) []*SearchResult {
	var searchResultList []*SearchResult
	for _, topic := range topics {
		searchResult := &SearchResult{
			Title: topic.Name,
			Arg:   fmt.Sprintf("%s/%s?project=%s", pubsubTopicEndpoint, topic.Name, gcpProject),
		}
		searchResultList = append(searchResultList, searchResult)
	}
	return searchResultList
}
