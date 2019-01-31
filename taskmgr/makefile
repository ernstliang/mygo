DIST := dist
EXECUTABLE := taskmgr-bin

GO ?= go
DEPLOY_ACCOUNT := industry
DEPLOY_IMAGE := $(EXECUTABLE)
GOFMT ?= gofmt "-s"

TARGETS ?= linux darwin windows
ARCHS ?= amd64
PACKAGES ?= $(shell $(GO) list ./... | grep -v /vendor/)
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
SOURCES ?= $(shell find . -name "*.go" -type f)
TAGS ?=
LDFLAGS ?= -X 'main.Version=$(VERSION)'
TMPDIR := $(shell mktemp -d 2>/dev/null || mktemp -d -t 'tempdir')
NODE_PROTOC_PLUGIN := $(shell which grpc_tools_node_protoc_plugin)
GOX := $(GOPATH)/bin/gox
MISSPELL := $(GOPATH)/bin/misspell

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
	SHA256SUM = sha256sum
else
	EXTLDFLAGS =
	SHA256SUM = shasum -a 256
endif

ifneq ($(DRONE_TAG),)
	VERSION ?= $(DRONE_TAG)
else
	VERSION ?= $(shell git describe --tags --always || git rev-parse --short HEAD)
endif

UNSUPPORT := To be continue

.PHONY: all
all: build

$(GOX):
	$(GO) get -u -v  github.com/mitchellh/gox

$(MISSPELL):
	$(GO) get -u -v github.com/client9/misspell/cmd/misspell

fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

vet:
	$(GO) vet $(PACKAGES)

module:
	$(GO) mod download

deps:
	$(GO) get -u -v github.com/campoy/embedmd

embedmd:
	embedmd -d *.md

errcheck:
	@hash errcheck > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u -v github.com/kisielk/errcheck; \
	fi
	errcheck $(PACKAGES)

lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u -v golang.org/x/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

unconvert:
	@hash unconvert > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u -v github.com/mdempsky/unconvert; \
	fi
	for PKG in $(PACKAGES); do sudo unconvert -v $$PKG || exit 1; done;

.PHONY: install
install: $(SOURCES)
	$(GO) install -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)'
	@echo "==> Installed ${EXECUTABLE} ${GOPATH}/bin/${EXECUTABLE}"

.PHONY: build
build: $(EXECUTABLE)
	$(foreach file,$(wildcard *.yaml),cp $(file) release/$(notdir $(file));)

$(EXECUTABLE): $(SOURCES)
	@echo $(EXECUTABLE)
	$(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o release/$@

.PHONY: misspell-check
misspell-check: $(MISSPELL)
	$(MISSPELL) -error $(GOFILES)

.PHONY: misspell
misspell: $(MISSPELL)
	$(MISSPELL) -w $(GOFILES)

.PHONY: coverage
coverage:
	@hash gocovmerge > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u -v github.com/wadey/gocovmerge; \
	fi
	gocovmerge $(shell find . -type f -name "coverage.out") > coverage.all;\

.PHONY: unit-test-coverage
unit-test-coverage:
	for PKG in $(PACKAGES); do $(GO) test -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || exit 1; done;

test:
	for PKG in $(PACKAGES); do $(GO) test -v $$PKG || exit 1; done;

config_test:
	$(GO) test -v -cover ./config/...

html:
	$(GO) tool cover -html=./coverage.all

release: release-dirs release-build release-copy release-check

release-dirs:
	mkdir -p $(DIST)/binaries $(DIST)/release

release-build: $(GOX)
	$(GOX) -os="$(TARGETS)" -arch="$(ARCHS)" -tags="$(TAGS)" -ldflags="$(EXTLDFLAGS)-s -w $(LDFLAGS)" -output="$(DIST)/binaries/$(EXECUTABLE)-$(VERSION)-{{.OS}}-{{.Arch}}"

release-copy:
	$(foreach file,$(wildcard $(DIST)/binaries/$(EXECUTABLE)-*),cp $(file) $(DIST)/release/$(notdir $(file));)
	$(foreach file,$(wildcard *.yaml),cp $(file) $(DIST)/release/$(notdir $(file));)

release-check:
	cd $(DIST)/release; $(foreach file,$(wildcard $(DIST)/binaries/$(EXECUTABLE)-*),$(SHA256SUM) $(notdir $(file)) > $(notdir $(file)).sha256;)

build_linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o release/linux/amd64/$(DEPLOY_IMAGE)

build_linux_i386:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o release/linux/i386/$(DEPLOY_IMAGE)

build_linux_arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o release/linux/arm64/$(DEPLOY_IMAGE)

build_linux_arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -a -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o release/linux/arm/$(DEPLOY_IMAGE)

build_linux_lambda:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags 'lambda' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o release/linux/lambda/$(DEPLOY_IMAGE)

docker_image:
	docker build -t $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE) -f Dockerfile .

docker_release: docker_image

docker_deploy:
ifeq ($(tag),)
	@echo "Usage: make $@ tag=<tag>"
	@exit 1
endif
	docker tag $(DEPLOY_ACCOUNT)/$(EXECUTABLE):latest $(DEPLOY_ACCOUNT)/$(EXECUTABLE):$(tag)
	docker push $(DEPLOY_ACCOUNT)/$(EXECUTABLE):$(tag)

clean:
	#$(GO) clean -modcache -x -i ./...
	$(GO) clean -x -i ./...
	find . -name "coverage.[ao][lu][lt]" -delete
	find . -name *.tar.gz -delete
	find . -name *.db -delete
	-rm -rf release dist

generate_proto:
	@echo $(UNSUPPORT)

version:
	@echo $(VERSION)

.PHONY: dev
dev:
	export server_config=development
	export MICRO_REGISTRY=mdns
	go run main.go
