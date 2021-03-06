GOOS = linux
GOARCH = amd64
GOBIN := $(GOPATH)/bin
GOLANGCILINT := $(GOBIN)/golangci-lint
GOSEC := $(GOBIN)/gosec
GHR := $(GOBIN)/ghr
RELEASEDIR ?= ../release/
CIRCLE_PROJECT_REPONAME ?= grace-log-parser
export GO111MODULE=on

.PHONY: build release clean test lint dependencies
build: test clean
	mkdir -p $(RELEASEDIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(RELEASEDIR)$(CIRCLE_PROJECT_REPONAME) -v
	zip -j $(RELEASEDIR)$(CIRCLE_PROJECT_REPONAME).zip $(RELEASEDIR)$(CIRCLE_PROJECT_REPONAME) email.html
	rm -f $(RELEASEDIR)$(CIRCLE_PROJECT_REPONAME)

release: build
	export CIRCLE_TAG=$(CIRCLE_TAG)
ifeq ($(strip $(GITHUB_TOKEN)),)
	@echo "GITHUB_TOKEN must be provided"
	@exit 0
else ifeq ($(strip $(CIRCLE_TAG)),)
	@echo "CIRCLE_TAG must be set"
	@exit 0
else ifeq ($(CIRCLE_BRANCH), "master")
	ghr -u $(CIRCLE_PROJECT_USERNAME) -r $(CIRCLE_PROJECT_REPONAME) -c $(CIRCLE_SHA1) -n $(CIRCLE_TAG) -delete $(CIRCLE_TAG) $(RELEASEDIR)
else
	ghr -prerelease -u $(CIRCLE_PROJECT_USERNAME) -r $(CIRCLE_PROJECT_REPONAME) -c $(CIRCLE_SHA1) -n $(CIRCLE_TAG) -delete $(CIRCLE_TAG) $(RELEASEDIR)
endif

clean:
	rm -rf $(RELEASEDIR)

test: lint
	go test -v ./...

lint: go.mod dependencies
	golangci-lint run ./...
	gosec ./...

go.mod:
ifeq (,$(wildcard go.mod))
	go mod init
endif

dependencies: $(GOLANGCILINT) $(GOSEC) $(GHR)

$(GOLANGCILINT):
	go get github.com/golangci/golangci-lint/cmd/golangci-lint

$(GOSEC):
	go get github.com/securego/gosec/cmd/gosec

$(GHR):
	go get github.com/tcnksm/ghr
