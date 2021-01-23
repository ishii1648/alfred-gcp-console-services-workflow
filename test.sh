#!/bin/bash

export alfred_workflow_bundleid="com.ishii1648.gcpconsoleservices"
export alfred_version=4.3
export alfred_workflow_version=1.0
export alfred_workflow_data="$HOME/Library/Application Support/$data_dir/Workflow Data/$alfred_workflow_bundleid"
export alfred_workflow_cache="$HOME/Library/Caches/$cache_dir/Workflow Data/$alfred_workflow_bundleid"

go test ./...