PKG = $(shell cat go.mod | grep "^module " | sed -e "s/module //g")
NAME = $(shell basename $(PKG))
VERSION = $(shell cat .version)
#VERSION = $(shell cat helmx.project.yml|grep version|awk -F : '{print $$2}'|tr -d " ")
COMMIT_SHA ?= $(shell git rev-parse --short HEAD)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= 0

GOBUILD=CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -ldflags "-X ${PKG}/version.Version=${VERSION}+sha.${COMMIT_SHA}"
PLATFORM := linux/amd64,linux/arm64
Github_UserName ?= 
Github_Token ?=

WORKSPACE ?= name

clean:
	rm -rf ./cmd/$(WORKSPACE)/out

upgrade:
	go get -u ./...

tidy:
	go mod tidy

build.srew.server: 
	cd ./cmd/$(WORKSPACE) && $(GOBUILD)


build: tidy
	$(MAKE) build.srew tar.srew GOOS=linux GOARCH=amd64
	$(MAKE) build.srew tar.srew GOOS=linux GOARCH=arm64


build.srew:
	cd ./cmd/$(WORKSPACE) && $(GOBUILD) -o ./out/srew-$(GOOS)-$(GOARCH)

tar.srew:
	cd ./cmd/$(WORKSPACE) && tar -czf ./out/srew-$(GOOS)-$(GOARCH).tar.gz -C ./out/ srew-$(GOOS)-$(GOARCH)

install: build.srew
	mv ./cmd/$(WORKSPACE)/out/srew-$(GOOS)-$(GOARCH) ${GOPATH}/bin/srew

docker.client:
	docker buildx build --push --progress plain --platform=${PLATFORM}	\
		--cache-from "type=local,src=/tmp/.buildx-cache" \
		--cache-to "type=local,dest=/tmp/.buildx-cache" \
		--file=./cmd/client/Dockerfile \
		--tag=bryantrh/srew:${VERSION}-${COMMIT_SHA} \
		--build-arg=Github_UserName=${Github_UserName}	\
		--build-arg=Github_Token=${Github_Token}	\
		.

docker.server:
	docker buildx build --push --progress plain --platform=${PLATFORM}	\
		--cache-from "type=local,src=/tmp/.buildx-cache" \
		--cache-to "type=local,dest=/tmp/.buildx-cache" \
		--file=./cmd/server/Dockerfile \
		--tag=bryantrh/srew-server:${VERSION}-${COMMIT_SHA} \
		--build-arg=Github_UserName=${Github_UserName}	\
		--build-arg=Github_Token=${Github_Token}	\
		.


gen-openapi:
	swag init --pd -d ./cmd/server -o ./cmd/server/docs

gen-client:
	swagger generate client -f ./cmd/server/docs/swagger.json -t ./cmd/client

gen-web:
	npx create-react-app web --template typescript


gen-web-client:
	restful-react import --file ./cmd/server/docs/swagger.json  --output ./cmd/web/src/client-bff.ts