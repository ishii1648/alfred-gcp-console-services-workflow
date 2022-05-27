package searchers

import (
	"context"
	"path/filepath"
	"runtime"
	"strings"

	aw "github.com/deanishe/awgo"
)

type Searcher interface {
	Search(ctx context.Context, wf *aw.Workflow, fullQuery string, gcpProject string, forceFetch bool) ([]*SearchResult, error)
}

type SearchResult struct {
	Title    string
	Subtitle string
	Arg      string
}

var gkeClustersSearcher = &GKEClustersSearcher{}
var pubSubTopicsSearcher = &PubSubTopicsSearcher{}
var pubSubSubscriptionsSearcher = &PubSubSubscriptionsSearcher{}
var gcsBrowserSearcher = &GcsBrowserSearcher{}
var cloudrunServicesSearcher = &CloudRunServicesSearcher{}
var projectSearcher = &ProjectSearcher{}
var cloudSQLSearcher = &CloudSQLSearcher{}

var SearchersByServiceId map[string]Searcher = map[string]Searcher{
	"gke_clusters":         gkeClustersSearcher,
	"pubsub_topics":        pubSubTopicsSearcher,
	"pubsub_subscriptions": pubSubSubscriptionsSearcher,
	"gcs_browser":          gcsBrowserSearcher,
	"cloudrun_services":    cloudrunServicesSearcher,
	"project":              projectSearcher,
	"sql_instances":        cloudSQLSearcher,
}

func getCurrentFilename() string {
	_, current_file, _, _ := runtime.Caller(1)
	baseFile := filepath.Base(current_file)
	return strings.TrimSuffix(baseFile, filepath.Ext(baseFile))
}
