// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package caching

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	aw "github.com/deanishe/awgo"
	"google.golang.org/api/container/v1"
)

type ContainerClusterListFetcher = func(ctx context.Context, gcpProject string) ([]*container.Cluster, error)

func LoadContainerClusterListFromCache(wf *aw.Workflow, ctx context.Context, cacheName string, fetcher ContainerClusterListFetcher, forceFetch bool, rawQuery string, gcpProject string) []*container.Cluster {
	cacheName += "_" + gcpProject
	results := []*container.Cluster{}
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
		return []*container.Cluster{}
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

type PubsubTopicListFetcher = func(ctx context.Context, gcpProject string) ([]*pubsub.Topic, error)

func LoadPubsubTopicListFromCache(wf *aw.Workflow, ctx context.Context, cacheName string, fetcher PubsubTopicListFetcher, forceFetch bool, rawQuery string, gcpProject string) []*pubsub.Topic {
	cacheName += "_" + gcpProject
	results := []*pubsub.Topic{}
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
		return []*pubsub.Topic{}
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
