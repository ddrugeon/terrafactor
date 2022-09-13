# syntax=docker/dockerfile:1

##
## STEP 1 - BUILD
##
FROM golang:1.19-alpine AS builder

LABEL maintainer="David Drugeon-Hamon <zebeurton@gmail.com>"

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . .

# Set necessary environment variables needed for our image and build the Create Go App CLI
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o terrafactor main.go

##
## STEP 2 - DEPLOY
##
FROM scratch

# Copy binary and config files from /build to root folder of scratch container.
COPY --from=builder ["/build/terrafactor", "/"]

# Set entry point.
ENTRYPOINT ["/terrafactor"]
