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

const pubsubSubEndpoint = "https://console.cloud.google.com/cloudpubsub/subscription/detail"

type PubSubSubscriptionsSearcher struct {
	c *pubsub.Client
}

func (s *PubSubSubscriptionsSearcher) Search(ctx context.Context, wf *aw.Workflow, fullQuery string, gcpProject string, forceFetch bool) ([]*SearchResult, error) {
	var searchResultList []*SearchResult
	c, err := pubsub.NewClient(ctx, gcpProject)
	if err != nil {
		return nil, err
	}
	s = &PubSubSubscriptionsSearcher{c: c}

	subscriptions := caching.LoadGcpPubsubSubscriptionListFromCache(wf, ctx, getCurrentFilename(), s.fetch, forceFetch, fullQuery, gcpProject)
	searchResultList = s.getSearchResultList(subscriptions, gcpProject)

	return searchResultList, nil
}

func (s *PubSubSubscriptionsSearcher) fetch(ctx context.Context, gcpProject string) ([]*gcp.PubsubSubscription, error) {
	var subscriptions []*gcp.PubsubSubscription
	it := s.c.Subscriptions(ctx)
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

func (s *PubSubSubscriptionsSearcher) getSearchResultList(subscriptions []*gcp.PubsubSubscription, gcpProject string) []*SearchResult {
	var searchResultList []*SearchResult
	for _, sub := range subscriptions {
		searchResult := &SearchResult{
			Title: sub.Name,
			Arg:   fmt.Sprintf("%s/%s?project=%s", pubsubSubEndpoint, sub.Name, gcpProject),
		}
		searchResultList = append(searchResultList, searchResult)
	}
	return searchResultList
}
