# FROM cgr.dev/chainguard/go:latest-dev AS builder
FROM golang:1.21.3-bullseye as builder

LABEL org.opencontainers.image.authors="raffael@nice.pink"
LABEL org.opencontainers.image.source="https://github.com/nice-pink/skupper-devops/blob/0f715e75a828a91ca1f4e0d3d61c6a8ac78dcb3a/Dockerfile"

WORKDIR /app

# get go module ready
COPY ./go.mod ./go.sum ./
RUN go mod download

# copy module code
COPY . .

# RUN mkdir -p bin
# RUN cd cmd/sitesync && go build -o ../../bin
RUN ./build_all

####################################################################################################

# FROM cgr.dev/chainguard/go:latest AS sitesync-runner
FROM golang:1.21.3-bullseye AS sitesync

WORKDIR /app

COPY --from=builder /app/bin/sitesync .

ENTRYPOINT [ "/app/sitesync" ]

####################################################################################################

FROM golang:1.21.3-bullseye AS autoheal

WORKDIR /app

COPY --from=builder /app/bin/autoheal .

ENTRYPOINT [ "/app/autoheal" ]

####################################################################################################

FROM golang:1.21.3-bullseye AS deploy

WORKDIR /app

COPY --from=builder /app/bin/deploy .

ENTRYPOINT [ "/app/deploy" ]
