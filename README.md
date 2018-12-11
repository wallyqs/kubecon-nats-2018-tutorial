### KubeCon 2018

#### Developing Production Ready Cloud Native NATS Applications

Download the source

```sh
git clone https://github.com/wallyqs/kubecon-nats-2018-tutorial.git
```

If you choose to just read along the just use the `complete` branch:

```sh
git checkout complete
```

#### Running the components

First get NATS Server and a NATS client:

```sh
go get -u github.com/nats-io/gnatsd
go get -u github.com/nats-io/go-nats
```

Starting the components:

```sh
gnatsd -DV -m 8222

# Starting the API Server
go run cmd/api-server/main.go

# Start the NYFT Agents
go run cmd/driver-agent/main.go

# Start the NYFT Service
go run cmd/rides-manager/main.go
```
