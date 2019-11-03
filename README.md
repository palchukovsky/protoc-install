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
* -ver: Protocol Buffers compiler version
* -out: output directory
### Example
```shell
    protoc-install -ver 3.10.1 -out bin
```