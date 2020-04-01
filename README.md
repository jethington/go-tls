# TLS Server and Client

## Description

This repository contains TLS server and client programs (server.go and client.go).  They connect on localhost to establish a TLS connection, then send and receive some example data.

This TLS setup has a slightly non-standard configuration:
* Forcing TLS 1.2
* Restricted set of cipher suites
* Client authentication

## References

The following links contain helpful documentation and/or example code that I used to create this project:
* https://golang.org/pkg/crypto/tls/
* https://gist.github.com/jim3ma/00523f865b8801390475c4e2049fe8c3
* https://gist.github.com/denji/12b3a568f092ab951456
* https://ericchiang.github.io/post/go-tls/
