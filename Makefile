PROJECT        = changelog-generator
GOLANG_VERSION = 1.11.2

BIN_DIR        = bin
PACKAGE        = $(shell go list)
PACKAGES       = $(shell go list -f '{{.Dir}}' ./... | grep -v /vendor/ )
GOBUILD        = CGO_ENABLED=0 go build
DOCKER_RUN     = docker run --rm -w /go/src/$(PACKAGE) -v $(CURDIR):/go/src/$(PACKAGE) golang:$(GOLANG_VERSION)

all: clean fmt lint test build docker-image

.PHONY: clean
vendors ?= 0
clean:
	rm -rf $(BIN_DIR)
ifeq ($(vendors),1)
	rm -rf vendor glide.lock
endif

.PHONY: fmt
fmt:
	for pkg in $(PACKAGES); do \
		gofmt -l -w -e $$pkg/*.go; \
	done

.PHONY: lint
lint:
	golangci-lint run

.PHONY: deps
deps:
	glide update; glide install

linux   ?= 1
darwin  ?= 0
windows ?= 0
.PHONY: build
build:
ifeq ($(linux),1)
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BIN_DIR)/$(PROJECT).linux $(PACKAGE)
endif
ifeq ($(darwin),1)
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BIN_DIR)/$(PROJECT).darwin $(PACKAGE)
endif
ifeq ($(windows),1)
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BIN_DIR)/$(PROJECT).exe $(PACKAGE)
endif

.PHONY: run
run:
	go run *.go -path=$(path)

.PHONY: test
test:
	go test -failfast -short -cover -v -timeout 10s -p=1 $(PACKAGES)

docker-%:
	$(DOCKER_RUN) make $* linux=$(linux) darwin=$(darwin) windows=$(windows)

docker-tag ?= local
.PHONY: docker-image
docker-image: docker-build
	docker build -t $(PROJECT):$(docker-tag) .
