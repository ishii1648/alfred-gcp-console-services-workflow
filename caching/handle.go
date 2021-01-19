package caching

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	aw "github.com/deanishe/awgo"
)

var jobName = "fetch"

func handleExpiredCache(wf *aw.Workflow, cacheName string, rawQuery string) error {
	maxCacheAgeSeconds := 180
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

	maxCacheAge := time.Duration(maxCacheAgeSeconds) * time.Second
	if wf.Cache.Expired(cacheName, maxCacheAge) {
		log.Printf("cache with key `%s` was expired (older than %d seconds) in %s", cacheName, maxCacheAgeSeconds, wf.CacheDir())
		wf.Rerun(0.5)
		if !wf.IsRunning(jobName) {
			cmd := exec.Command(os.Args[0], "-query="+rawQuery+"", "-fetch")
			log.Printf("running `%s` in background as job `%s` ...", cmd, jobName)
			if err := wf.RunInBackground(jobName, cmd); err != nil {
				panic(err)
			}
		} else {
			log.Printf("background job `%s` already running", jobName)
		}
	}

	return nil
}
