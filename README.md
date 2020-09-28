# File Service

This is the File Service

Generated with

```
micro new gitlab.com/bimoyong/go-file --namespace=go --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.srv.file
- Type: srv
- Alias: file

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend consul.

```
# install consul
brew install consul

# run consul
consul agent -dev
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./bin/srv.file
```

Build a docker image
```
make docker
```