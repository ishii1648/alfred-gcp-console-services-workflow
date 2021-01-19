package searchers

import (
	"context"
	"fmt"

	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/caching"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/gcp"
	"google.golang.org/api/container/v1"
)

type GKEClustersSearcher struct{}

func (s *GKEClustersSearcher) Search(ctx context.Context, wf *aw.Workflow, fullQuery string, gcpProject string, gcpService gcp.GcpService, forceFetch bool) error {
	cacheName := getCurrentFilename()
	clusters := caching.LoadContainerClusterListFromCache(wf, ctx, cacheName, s.fetch, forceFetch, fullQuery, gcpProject)

	for _, cluster := range clusters {
		s.addToWorkflow(wf, cluster, gcpService, gcpProject)
	}
	return nil
}

func (s *GKEClustersSearcher) fetch(ctx context.Context, gcpProject string) ([]*container.Cluster, error) {
	containerService, err := container.NewService(ctx)
	if err != nil {
		return nil, err
	}

	clusters, err := container.NewProjectsLocationsClustersService(containerService).List(fmt.Sprintf("projects/%s/locations/-", gcpProject)).Do()
	if err != nil {
		return nil, err
	}

	return clusters.Clusters, nil
}

func (s *GKEClustersSearcher) addToWorkflow(wf *aw.Workflow, cluster *container.Cluster, gcpService gcp.GcpService, gcpProject string) {
	wf.NewItem(cluster.Name).
		Valid(true).
		Var("action", "open-url").
		Subtitle(fmt.Sprintf("%s %s", s.getStatusEmoji(cluster.CurrentNodeCount), cluster.Location)).
		Arg(fmt.Sprintf("https://console.cloud.google.com/kubernetes/clusters/details/%s/%s?project=%s", cluster.Location, cluster.Name, gcpProject)).
		Icon(&aw.Icon{Value: gcpService.GetIcon()})
}

func (s *GKEClustersSearcher) getStatusEmoji(clusterSize int64) string {
	if clusterSize == 0 {
		return "🔴"
	} else {
		return "🟢"
	}
}
