# FROM cgr.dev/chainguard/go:latest-dev AS builder
FROM golang:1.21.3-bullseye as builder

LABEL org.opencontainers.image.authors="raffael@nice.pink"
LABEL org.opencontainers.image.source="https://github.com/radiosphere/rs-ops/blob/5243f8d4d53809f90fdc83d47a9dc71611743c77/resources/services/Dockerfile"

WORKDIR /app

# get go module ready
COPY ./go.mod ./go.sum ./
RUN go mod download

# copy module code
COPY . .

RUN mkdir -p bin
RUN cd cmd/sitesync && go build -o ../../bin

####################################################################################################

# FROM builder as sitesync-builder

# RUN go build cmd/sitesync

# FROM cgr.dev/chainguard/go:latest AS sitesync-runner
FROM golang:1.21.3-bullseye AS sitesync-runner

WORKDIR /app

COPY --from=builder /app/bin/sitesync .

ENTRYPOINT [ "/app/sitesync" ]
