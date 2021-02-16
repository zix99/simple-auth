TAG ?= head
COMMIT_SHA ?= devel

all: build

build: generate swagger server

server:
	go build -ldflags "-X main.version=${TAG} -X main.buildSha=${COMMIT_SHA}" -tags box,swagger,prometheus -o bin/simple-auth-server simple-auth/cmd/server

cli:
	go build -ldflags "-X main.version=${TAG} -X main.buildSha=${COMMIT_SHA}" -tags boxconfig -o bin/simple-auth-cli simple-auth/cmd/cli

rundev:
	go run -tags swagger,prometheus simple-auth/cmd/server --include=simpleauth.yml

generate:
	go generate ./...

swagger:
	go run github.com/swaggo/swag/cmd/swag init -o pkg/swagdocs -g pkg/routes/api/api.go

# TESTING

integrationtest:
	./tests/quicktest.sh

unittest:
	go test -race ./...

vet:
	go vet $$(go list ./... | grep -v "vendor")

staticcheck:
	go run honnef.co/go/tools/cmd/staticcheck ./...


test: unittest vet

check: test staticcheck integrationtest

# CLEANING

clean:
	rm -rf dist/
	rm -f pkg/box/*.gen.go
	rm -rf pkg/swagdocs
