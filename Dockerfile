# Build node app
FROM node:12-slim AS nodebuild
WORKDIR /opt/simple-auth
COPY package*.json ./
RUN npm ci
COPY webpack.config.js .
COPY vue vue
RUN npm run build

# Build go app
FROM golang:1.14-alpine AS gobuild
RUN apk add build-base
WORKDIR /opt/simple-auth
COPY --from=nodebuild /opt/simple-auth/dist dist
COPY . .
RUN go generate ./...
RUN go build -tags box -o simple-auth-server simple-auth/cmd/server

# Final image
FROM alpine:latest
WORKDIR /opt/simple-auth
COPY --from=gobuild /opt/simple-auth/simple-auth-server .

VOLUME /var/lib/simple-auth
ENV SA_PRODUCTION=true \
    SA_WEB_HOST="0.0.0.0:80" \
    SA_DB_URL="/var/lib/simple-auth/simpleauth.db"

CMD ["./simple-auth-server"]
