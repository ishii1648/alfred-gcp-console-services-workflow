PLIST=info.plist
EXEC_BIN=alfred-gcp-console-services-workflow
DIST_FILE=alfred-gcp-console-services-workflow.alfredworkflow
GO_SRCS=$(shell find -f . \( -name \*.go \))

.PHONY: build
build:
	go build -o $(EXEC_BIN)
	zip -r $(DIST_FILE) images console-services.yml $(EXEC_BIN) $(PLIST) icon.png

.PHONY: test
test:
	bash ./test.sh