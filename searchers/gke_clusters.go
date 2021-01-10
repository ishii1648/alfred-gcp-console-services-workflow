package searchers

import (
	"context"
	"fmt"

	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/gcp"
	"google.golang.org/api/container/v1"
)

type GKEClustersSearcher struct {
	gcpProject string
	gcpService gcp.GcpService
}

func (s *GKEClustersSearcher) Search(ctx context.Context, wf *aw.Workflow, gcpProject string, gcpService gcp.GcpService) error {
	s = &GKEClustersSearcher{
		gcpProject: gcpProject,
		gcpService: gcpService,
	}

	clusters, err := s.fetch(ctx)
	if err != nil {
		return err
	}

	for _, cluster := range clusters {
		s.addToWorkflow(wf, cluster)
	}
	return nil
}

func (s *GKEClustersSearcher) fetch(ctx context.Context) ([]*container.Cluster, error) {
	containerService, err := container.NewService(ctx)
	if err != nil {
		return nil, err
	}

	clusters, err := container.NewProjectsLocationsClustersService(containerService).List(fmt.Sprintf("projects/%s/locations/-", s.gcpProject)).Do()
	if err != nil {
		return nil, err
	}

	return clusters.Clusters, nil
}

func (s *GKEClustersSearcher) addToWorkflow(wf *aw.Workflow, cluster *container.Cluster) {
	wf.NewItem(cluster.Name).
		Valid(true).
		Var("action", "open-url").
		Subtitle(fmt.Sprintf("%s %s", s.getStatusEmoji(cluster.CurrentNodeCount), cluster.Location)).
		Arg(fmt.Sprintf("https://console.cloud.google.com/kubernetes/clusters/details/%s/%s?project=%s", cluster.Location, cluster.Name, s.gcpProject)).
		Icon(&aw.Icon{Value: s.gcpService.GetIcon()})
}

func (s *GKEClustersSearcher) getStatusEmoji(clusterSize int64) string {
	if clusterSize == 0 {
		return "🔴"
	} else {
		return "🟢"
	}
}
