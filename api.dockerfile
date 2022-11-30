ARG INTERNAL_REGISTRY=""
FROM ${INTERNAL_REGISTRY}public.ecr.aws/docker/library/golang:1.18 as build

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /our-code

COPY / ./
RUN ls
RUN go build ./cmd/main/
RUN ls
RUN go mod vendor


RUN go run ./cmd/main
