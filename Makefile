UNAME := $(shell uname)
ifneq (,$(findstring MINGW,$(UNAME)))
#Gopath is not saved across sessions, probably existing Windows env vars, override them
export GOPATH := $(HOME)/go
GOPATH1 := $(GOPATH)
export PATH := $(PATH):$(GOPATH)/bin
else
export GOPATH := $(shell go env GOPATH)
GOPATH1 := $(firstword $(subst :, ,$(GOPATH)))
endif
SRCPATH     := $(shell pwd)
ARCH        := $(shell ./zarf/scripts/archtype.sh)
OS_TYPE     := $(shell ./zarf/scripts/ostype.sh)

#Give execute permission for the scripts
$(shell chmod -R u+x ./zarf/scripts/)

# If build number already set, use it - to ensure same build number across multiple platforms being built
BUILDNUMBER      ?= $(shell ./zarf/scripts/compute_build_number.sh)
FULLBUILDNUMBER  ?= $(shell ./zarf/scripts/compute_build_number.sh -f)
COMMITHASH       := $(shell ./zarf/scripts/compute_build_commit.sh)
BUILDBRANCH      := $(shell ./zarf/scripts/compute_branch.sh)
CHANNEL          ?= $(shell ./zarf/scripts/compute_branch_channel.sh $(BUILDBRANCH))

export GOCACHE=$(SRCPATH)/build/go-cache
export GOTESTCOMMAND=go test

ifneq (, $(findstring MINGW,$(UNAME)))
export GOBUILDMODE := -buildmode=exe
endif

GOTRIMPATH	:= $(shell GOPATH=$(GOPATH) && go help build | grep -q .-trimpath && echo -trimpath)

GOLDFLAGS  := -X crossjoin.com/gorxestra/util/conf.BuildNumber=$(BUILDNUMBER) \
		 -X crossjoin.com/gorxestra/util/conf.CommitHash=$(COMMITHASH) \
		 -X crossjoin.com/gorxestra/util/conf.Branch=$(BUILDBRANCH) \
		 -X crossjoin.com/gorxestra/util/conf.Channel=$(CHANNEL)

DOCKER_FILE := ./zarf/docker/Dockerfile
IMAGE_NAME := gorxestra
IMAGE_TAG := $(CHANNEL)

default: build 

generate:
	go generate ./...

build:
	go build -o ./build/ $(GOTRIMPATH) $(GOBUILDMODE) -ldflags="$(GOLDFLAGS)" ./...

docker:
	docker build \
		-f $(DOCKER_FILE) \
		-t $(IMAGE_NAME):$(IMAGE_TAG) \
		--build-arg VERSION=$(FULLBUILDNUMBER) \
		--build-arg COMMIT_SHA=$(COMMITHASH) \
		. 

docker-test: docker-clean 
	mkdir -p ./tmp/musician
	mkdir -p ./tmp/conductor
	chown -R 1000:1000 ./tmp
	docker build \
		-f $(DOCKER_FILE) \
		-t $(IMAGE_NAME):test \
		--build-arg VERSION=$(FULLBUILDNUMBER) \
		--build-arg COMMIT_SHA=$(COMMITHASH) \
		. 
	docker compose -f ./zarf/docker/docker-compose.yml up -d
	docker compose -f ./zarf/docker/docker-compose.yml logs -f

run: build
	rm -rf tmp
	mkdir -p ./tmp/conductor
	cp files/conductor/* tmp/conductor

	ENV_REST_ENDPOINTADDRESS=0.0.0.0:8080 \
	./build/conductor -d=tmp/conductor > tmp/conductor/logs.txt &

	./launch.sh

stop:
	pkill -f "./build/" && \
	rm -rf tmp

docker-clean:
	docker compose -f ./zarf/docker/docker-compose.yml down

clean: docker-clean	
	go clean -i ./...
	rm -rf ./build
	rm -rf ./tmp

deps:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

tidy:
	go mod tidy
	go mod vendor

.PHONY: default build clean
