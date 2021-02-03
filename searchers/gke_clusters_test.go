package searchers

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	aw "github.com/deanishe/awgo"
	gax "github.com/googleapis/gax-go/v2"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/caching"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

const testGcpProject = "sample-gcp-project"

type fakeGKEClustersSearcher struct {
	GKECluster
	FakeListClusters func(ctx context.Context, req *containerpb.ListClustersRequest, opts ...gax.CallOption) (*containerpb.ListClustersResponse, error)
}

func (s *fakeGKEClustersSearcher) ListClusters(ctx context.Context, req *containerpb.ListClustersRequest, opts ...gax.CallOption) (*containerpb.ListClustersResponse, error) {
	return s.FakeListClusters(ctx, req)
}

type GKEClustersTest struct {
	name       string
	forceFetch bool
	clusters   []*containerpb.Cluster
	wants      []*SearchResult
}

var gkeClustersTests []GKEClustersTest = []GKEClustersTest{
	{
		name:       "test1",
		forceFetch: true,
		clusters: []*containerpb.Cluster{
			{
				Name:             "test-cluster01",
				CurrentNodeCount: 1,
				Location:         "asia-east1",
			},
		},
		wants: []*SearchResult{
			{
				Title:    "test-cluster01",
				Subtitle: "🟢 asia-east1",
				Arg:      fmt.Sprintf("%s/%s/%s?project=%s", gkeClusterEndpoint, "asia-east1", "test-cluster01", testGcpProject),
			},
		},
	},
}

func TestSearch(t *testing.T) {
	var searchResultList []*SearchResult
	wf := aw.New()

	for _, gt := range gkeClustersTests {
		t.Run(gt.name, func(t *testing.T) {
			fakeClient := &fakeGKEClustersSearcher{
				FakeListClusters: func(ctx context.Context, req *containerpb.ListClustersRequest, opts ...gax.CallOption) (*containerpb.ListClustersResponse, error) {
					return &containerpb.ListClustersResponse{Clusters: gt.clusters}, nil
				},
			}

			s := &GKEClustersSearcher{gke: fakeClient}
			ctx := context.Background()

			clusters := caching.LoadContainerpbClusterListFromCache(wf, ctx, getCurrentFilename(), s.fetch, gt.forceFetch, "", testGcpProject)
			searchResultList = s.getSearchResultList(clusters, testGcpProject)

			for i, searchResult := range searchResultList {
				if !reflect.DeepEqual(gt.wants[i], searchResult) {
					t.Errorf("want = %v, got = %v", gt.wants[i], searchResult)
				}
			}
		})
	}
}
