package searchers

import (
	"context"
	"fmt"

	container "cloud.google.com/go/container/apiv1"
	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/caching"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/gcp"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

type GKEClustersSearcher struct{}

func (s *GKEClustersSearcher) Search(ctx context.Context, wf *aw.Workflow, fullQuery string, gcpProject string, gcpService gcp.GcpService, forceFetch bool) error {
	cacheName := getCurrentFilename()
	clusters := caching.LoadContainerpbClusterListFromCache(wf, ctx, cacheName, s.fetch, forceFetch, fullQuery, gcpProject)

	for _, cluster := range clusters {
		s.addToWorkflow(wf, cluster, gcpService, gcpProject)
	}
	return nil
}

func (s *GKEClustersSearcher) fetch(ctx context.Context, gcpProject string) ([]*containerpb.Cluster, error) {
	c, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		return nil, err
	}

	req := &containerpb.ListClustersRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", gcpProject, "-"),
	}
	resp, err := c.ListClusters(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Clusters, nil
}

func (s *GKEClustersSearcher) addToWorkflow(wf *aw.Workflow, cluster *containerpb.Cluster, gcpService gcp.GcpService, gcpProject string) {
	wf.NewItem(cluster.Name).
		Valid(true).
		Var("action", "open-url").
		Subtitle(fmt.Sprintf("%s %s", s.getStatusEmoji(cluster.CurrentNodeCount), cluster.Location)).
		Arg(fmt.Sprintf("https://console.cloud.google.com/kubernetes/clusters/details/%s/%s?project=%s", cluster.Location, cluster.Name, gcpProject)).
		Icon(&aw.Icon{Value: gcpService.GetIcon()})
}

func (s *GKEClustersSearcher) getStatusEmoji(clusterSize int32) string {
	if clusterSize == 0 {
		return "🔴"
	} else {
		return "🟢"
	}
}
