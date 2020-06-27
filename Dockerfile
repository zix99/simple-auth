# Build node app
FROM node:12-slim AS nodebuild
WORKDIR /opt/simple-auth
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Build go app
FROM golang:1.14-alpine AS gobuild
RUN apk add build-base
WORKDIR /opt/simple-auth
COPY . .
RUN go build -o simple-auth-server simple-auth/cmd/server

# Final image
FROM alpine:latest
WORKDIR /opt/simple-auth
COPY static static
COPY templates templates
COPY simpleauth.default.yml .
COPY --from=nodebuild /opt/simple-auth/dist .
COPY --from=gobuild /opt/simple-auth/simple-auth-server .
CMD ["./simple-auth-server"]