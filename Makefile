PLIST=info.plist
EXEC_BIN=alfred-gcp-console-services-workflow
DIST_FILE=alfred-gcp-console-services-workflow.alfredworkflow
GO_SRCS=$(shell find -f . \( -name \*.go \))

.PHONY: build
build:
	go build -o $(EXEC_BIN)
	cp console-services.yml $(EXEC_BIN) ~/Library/Application\ Support/Alfred/Alfred.alfredpreferences/workflows/user.workflow.5CDC68D2-492C-4943-BFE1-1F91D789954A/
	cp -ar images ~/Library/Application\ Support/Alfred/Alfred.alfredpreferences/workflows/user.workflow.5CDC68D2-492C-4943-BFE1-1F91D789954A/