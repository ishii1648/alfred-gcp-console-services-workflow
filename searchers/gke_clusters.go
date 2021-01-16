package searchers

import (
	"context"
	"fmt"
	"log"
	"os"
	// "os/exec"
	"strconv"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/ishii1648/alfred-gcp-console-services-workflow/gcp"
	"google.golang.org/api/container/v1"
)

type GKEClustersSearcher struct {
	gcpProject string
	gcpService gcp.GcpService
}

func (s *GKEClustersSearcher) Search(ctx context.Context, wf *aw.Workflow, fullQuery string, gcpProject string, gcpService gcp.GcpService, forceFetch bool) error {
	s = &GKEClustersSearcher{
		gcpProject: gcpProject,
		gcpService: gcpService,
	}

	var clusters []*container.Cluster

	cacheName := fmt.Sprintf("cache_gkecluster_%s.json", gcpProject)
	log.Printf("forceFetch : %v", forceFetch)
	if forceFetch {
		log.Printf("fetching from gcp ...")

		clusters, err := s.fetch(ctx)
		if err != nil {
			return err
		}

		log.Printf("storing %d results with cache key `%s` to %s ...", len(clusters), cacheName, wf.CacheDir())
		if err := wf.Cache.StoreJSON(cacheName, clusters); err != nil {
			panic(err)
		}
	}

	maxCacheAgeSeconds := 15
	m := os.Getenv("ALFRED_GCP_CONSOLE_SERVICES_WORKFLOW_MAX_CACHE_AGE_SECONDS")
	if m != "" {
		converted, err := strconv.Atoi(m)
		if err != nil {
			panic(err)
		}
		if converted != 0 {
			log.Printf("using custom max cache age of %v seconds", converted)
			maxCacheAgeSeconds = converted
		}
	}

	// jobName := "fetch"
	maxCacheAge := time.Duration(maxCacheAgeSeconds) * time.Second
	if wf.Cache.Expired(cacheName, maxCacheAge) {
		log.Printf("cache with key `%s` was expired (older than %d seconds) in %s", cacheName, maxCacheAgeSeconds, wf.CacheDir())
		// wf.Rerun(0.5)
		// if !wf.IsRunning(jobName) {
		// 	cmd := exec.Command(os.Args[0], "-query="+fullQuery+"", "-fetch")
		// 	if err := wf.RunInBackground(jobName, cmd); err != nil {
		// 		log.Printf("failed to run background job : %v", err)
		// 		return err
		// 	}
		// 	log.Printf("running `%s` in background as job `%s` ...", cmd, jobName)
		// } else {
		// 	log.Printf("background job `%s` already running", jobName)
		// }
	}

	if wf.Cache.Exists(cacheName) {
		log.Printf("using cache with key `%s` in %s ...", cacheName, wf.CacheDir())
		if err := wf.Cache.LoadJSON(cacheName, &clusters); err != nil {
			panic(err)
		}
	} else {
		log.Printf("cache with key `%s` did not exist in %s ...", cacheName, wf.CacheDir())
		wf.NewItem("Fetching ...").
			Icon(aw.IconInfo)
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
