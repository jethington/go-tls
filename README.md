# TLS Server and Client

## Description

This repository contains TLS server and client programs (server.go and client.go).  They connect on localhost to establish a TLS connection, then send and receive some example data.

This TLS setup has a slightly non-standard configuration:
* Forcing TLS 1.2
* Restricted set of cipher suites
* Client authentication
