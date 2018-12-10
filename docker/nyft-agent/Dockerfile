FROM golang:1.11-alpine3.8 AS builder

WORKDIR $GOPATH/src/github.com/wallyqs/kubecon-nats-2018-tutorial/

MAINTAINER Waldemar Quevedo <wally@synadia.com>

RUN apk add --update git

RUN go get -u github.com/nats-io/go-nats

COPY . .

RUN CGO_ENABLED=0 go build -o /nyft-agent -v -a -tags netgo -installsuffix netgo cmd/nyft-agent/main.go

FROM alpine:3.8

RUN apk add --update ca-certificates

COPY --from=builder /nyft-agent /bin/nyft-agent

ENTRYPOINT ["/bin/nyft-agent"]
