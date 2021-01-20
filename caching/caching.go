package caching

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/cheekybits/genny/generic"
	aw "github.com/deanishe/awgo"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "Entity=*containerpb.Cluster,*storage.BucketAttrs,*gcp.PubsubSubscription,*gcp.PubsubTopic"
type Entity = generic.Type

type EntityListFetcher = func(ctx context.Context, gcpProject string) ([]Entity, error)

func LoadEntityListFromCache(wf *aw.Workflow, ctx context.Context, cacheName string, fetcher EntityListFetcher, forceFetch bool, rawQuery string, gcpProject string) []Entity {
	cacheName += "_" + gcpProject
	results := []Entity{}
	lastFetchErrPath := wf.CacheDir() + "/last-fetch-err.txt"

	if forceFetch {
		log.Printf("fetching from gcp ...")
		results, err := fetcher(ctx, gcpProject)

		if err != nil {
			log.Printf("fetch error occurred. writing to %s ...", lastFetchErrPath)
			ioutil.WriteFile(lastFetchErrPath, []byte(err.Error()), 0600)
			panic(err)
		} else {
			os.Remove(lastFetchErrPath)
		}

		log.Printf("storing %d results with cache key `%s` to %s ...", len(results), cacheName, wf.CacheDir())
		if err := wf.Cache.StoreJSON(cacheName, results); err != nil {
			panic(err)
		}

		return results
	}

	err := handleExpiredCache(wf, cacheName, lastFetchErrPath, rawQuery)
	if err != nil {
		return []Entity{}
	}

	if wf.Cache.Exists(cacheName) {
		log.Printf("using cache with key `%s` in %s ...", cacheName, wf.CacheDir())
		if err := wf.Cache.LoadJSON(cacheName, &results); err != nil {
			panic(err)
		}
	} else {
		log.Printf("cache with key `%s` did not exist in %s ...", cacheName, wf.CacheDir())
		wf.NewItem("Fetching ...").
			Icon(aw.IconInfo)
	}

	return results
}
