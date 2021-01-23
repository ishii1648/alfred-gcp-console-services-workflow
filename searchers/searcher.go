package searchers

import (
	"context"
	"path/filepath"
	"runtime"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/gcp"
)

type Searcher interface {
	Search(ctx context.Context, wf *aw.Workflow, query string, gcpProject string, gcpService gcp.GcpService, forceFetch bool) error
}

var gkeClustersSearcher = &GKEClustersSearcher{}
var pubSubTopicsSearcher = &PubSubTopicsSearcher{}
var pubSubSubscriptionsSearcher = &PubSubSubscriptionsSearcher{}
var gcsBrowserSearcher = &GcsBrowserSearcher{}
var cloudrunServicesSearcher = &CloudRunServicesSearcher{}

var SearchersByServiceId map[string]Searcher = map[string]Searcher{
	"gke_clusters":         gkeClustersSearcher,
	"pubsub_topics":        pubSubTopicsSearcher,
	"pubsub_subscriptions": pubSubSubscriptionsSearcher,
	"gcs_browser":          gcsBrowserSearcher,
	"cloudrun_services":    cloudrunServicesSearcher,
}

func getCurrentFilename() string {
	_, current_file, _, _ := runtime.Caller(1)
	baseFile := filepath.Base(current_file)
	return strings.TrimSuffix(baseFile, filepath.Ext(baseFile))
}
