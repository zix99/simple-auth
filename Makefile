TAG ?= head
COMMIT_SHA ?= devel
OUT ?= bin/

all: vue build

vue:
	npm run build

# GO Binary (May depend on npm commands first)

build: generate swagger bin

bin: server cli

server:
	go build -ldflags "-X main.version=${TAG} -X main.buildSha=${COMMIT_SHA}" -tags box,swagger,prometheus -o ${OUT}simple-auth-server simple-auth/cmd/server

cli:
	go build -ldflags "-X main.version=${TAG} -X main.buildSha=${COMMIT_SHA}" -tags boxconfig -o ${OUT}simple-auth-cli simple-auth/cmd/cli

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
