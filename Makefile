all: build

build: generate swagger
	go build -ldflags "-X main.version=${DRONE_TAG:-head} -X main.buildSha=${DRONE_COMMIT_SHA}" -tags box,swagger,prometheus -o bin/simple-auth-server simple-auth/cmd/server

rundev:
	go run -tags swagger,prometheus simple-auth/cmd/server --include=simpleauth.yml

generate:
	go generate ./...

swagger:
	go run github.com/swaggo/swag/cmd/swag init -o pkg/swagdocs -g pkg/routes/api/api.go

clean:
	rm -rf dist/
	rm -f pkg/box/*.gen.go
	rm -rf pkg/swagdocs
