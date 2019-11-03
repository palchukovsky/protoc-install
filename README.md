[![Go Report Card](https://goreportcard.com/badge/github.com/palchukovsky/protoc-install)](https://goreportcard.com/report/github.com/palchukovsky/protoc-install)

# protoc-install
Downloads and unpacks Protocol Buffers compiler protoc.

## To install
[Go](https://golang.org/dl/) has to be installed.
```shell
    go get github.com/palchukovsky/protoc-install
```

## To run
### Arguments
* -type: "cli" for protoc or "grpc-web" for gRPC Web protoc plugin
* -ver: version to install
* -out: output directory
### Example
```shell
    protoc-install -type cli -ver 3.10.1 -out bin
```