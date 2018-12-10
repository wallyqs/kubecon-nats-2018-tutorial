FROM golang:1.11-alpine3.8 AS builder

WORKDIR $GOPATH/src/github.com/wallyqs/kubecon-nats-2018-tutorial/

MAINTAINER Waldemar Quevedo <wally@synadia.com>

RUN apk add --update git

RUN go get -u github.com/nats-io/go-nats

COPY . .

RUN CGO_ENABLED=0 go build -o /rides-manager -v -a -tags netgo -installsuffix netgo cmd/rides-manager/main.go

FROM alpine:3.8

RUN apk add --update ca-certificates

COPY --from=builder /rides-manager /bin/rides-manager

ENTRYPOINT ["/bin/rides-manager"]
