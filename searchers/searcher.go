package searchers

import (
	"context"
	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/gcp"
)

type Searcher interface {
	Search(ctx context.Context, wf *aw.Workflow, gcpProject string, gcpService gcp.GcpService) error
}

var gkeClustersSearcher = &GKEClustersSearcher{}

var SearchersByServiceId map[string]Searcher = map[string]Searcher{
	"gke_clusters": gkeClustersSearcher,
}

// func GetStateEmoji(state string) string {
// 	switch state {
// 	case ec2.InstanceStateNameRunning:
// 		return "🟢"
// 	case ec2.InstanceStateNameShuttingDown:
// 		return "🟡"
// 	case ec2.InstanceStateNameStopping:
// 		return "🟡"
// 	case ec2.InstanceStateNameStopped:
// 		return "🔴"
// 	case ec2.InstanceStateNameTerminated:
// 		return "🔴"
// 	case ec2.InstanceStateNamePending:
// 		return "⚪️"
// 	}

// 	return "❔"
// }
