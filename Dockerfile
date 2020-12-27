# Build node app
FROM node:12-slim AS nodebuild
WORKDIR /opt/simple-auth
COPY package*.json ./
RUN npm ci
COPY webpack.config.js .
COPY vue vue
RUN npm run build

# Build go app
FROM golang:1.15-alpine AS gobuild
RUN apk add build-base
WORKDIR /opt/simple-auth
COPY go.* ./
RUN go mod download
COPY . .
COPY --from=nodebuild /opt/simple-auth/dist dist
RUN go generate ./...
RUN go run github.com/swaggo/swag/cmd/swag init -o pkg/swagdocs -g pkg/routes/api/api.go
RUN go build -tags box,prometheus,swagger -o simple-auth-server simple-auth/cmd/server
RUN go build -tags boxconfig -o simple-auth-cli simple-auth/cmd/cli

# Final image
FROM alpine:latest
WORKDIR /opt/simple-auth
COPY --from=gobuild /opt/simple-auth/simple-auth-server .
COPY --from=gobuild /opt/simple-auth/simple-auth-cli .

VOLUME /var/lib/simple-auth
ENV SA_PRODUCTION=true \
    SA_WEB_HOST="0.0.0.0:80" \
    SA_DB_URL="/var/lib/simple-auth/simpleauth.db"

EXPOSE 80
ENTRYPOINT ["./simple-auth-server"]
CMD []
