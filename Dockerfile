#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git

# CGO has to be disabled for alpine
ENV CGO_ENABLED=0

WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/forklift /forklift
ENTRYPOINT ./forklift

LABEL Name=forklift Version=0.0.1
LABEL org.opencontainers.image.source https://github.com/dacort/forklift