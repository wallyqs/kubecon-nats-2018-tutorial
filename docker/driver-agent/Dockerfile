FROM golang:1.11-alpine3.8 AS builder

WORKDIR $GOPATH/src/github.com/wallyqs/kubecon-nats-2018-tutorial/

MAINTAINER Waldemar Quevedo <wally@synadia.com>

RUN apk add --update git

RUN go get -u github.com/nats-io/go-nats

COPY . .

RUN CGO_ENABLED=0 go build -o /driver-agent -v -a -tags netgo -installsuffix netgo cmd/driver-agent/main.go

FROM alpine:3.8

RUN apk add --update ca-certificates

COPY --from=builder /driver-agent /bin/driver-agent

ENTRYPOINT ["/bin/driver-agent"]
