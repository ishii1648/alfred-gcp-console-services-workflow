package caching

import (
	"context"
	"log"

	"github.com/cheekybits/genny/generic"
	aw "github.com/deanishe/awgo"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "Entity=*container.Cluster"
type Entity = generic.Type

type EntityArrayFetcher = func(ctx context.Context, gcpProject string) ([]Entity, error)

func LoadEntityArrayFromCache(wf *aw.Workflow, ctx context.Context, cacheName string, fetcher EntityArrayFetcher, forceFetch bool, rawQuery string, gcpProject string) []Entity {
	results := []Entity{}
	if forceFetch {
		log.Printf("fetching from gcp ...")

		results, err := fetcher(ctx, gcpProject)
		if err != nil {
			log.Printf("failed to fetcher : %v", err)
			panic(err)
		}

		log.Printf("storing %d results with cache key `%s` to %s ...", len(results), cacheName, wf.CacheDir())
		if err := wf.Cache.StoreJSON(cacheName, results); err != nil {
			panic(err)
		}

		return results
	}

	err := handleExpiredCache(wf, cacheName, rawQuery)
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
