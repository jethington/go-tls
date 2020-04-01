package main

import (
    "crypto/tls"
    "log"
    "net"
    "bufio"
    "crypto/x509"
    "io/ioutil"
)

func main() {
    cert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
    if err != nil {
        log.Fatalf("SERVER: tls.LoadX509KeyPair: %s", err) // TODO investigate log.Fatalf vs log.Println + return
    }

    dat, err := ioutil.ReadFile("certs/root.crt")
    if err != nil {
        log.Println(err)
        return
    }
    roots := x509.NewCertPool()
    ok := roots.AppendCertsFromPEM(dat)
    if !ok {
        log.Println("failed to parse root certificate")
        return
    }

    cfg := tls.Config{
        MinVersion:                 tls.VersionTLS12,
        MaxVersion:                 tls.VersionTLS12,
        CurvePreferences:           []tls.CurveID{tls.CurveP256},
        PreferServerCipherSuites:   true,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
            tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
        },
        Certificates:               []tls.Certificate{cert},
        ClientCAs:                  roots,
        SessionTicketsDisabled:     true,
        ClientAuth:                 tls.RequireAndVerifyClientCert,
    }

    ln, err := tls.Listen("tcp", "127.0.0.1:8883", &cfg)
    if err != nil {
        log.Fatalf("SERVER: tls.Listen: %s", err)
    }
    defer ln.Close()
    log.Print("server: listening")

    count := 0
    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Println(err)
            continue
        }
        go handleConnection(conn, count)
        count += 1
    }
}

func handleConnection(conn net.Conn, count int) {
    defer conn.Close()
    reader := bufio.NewReader(conn)
    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            log.Println(err)
            return
        }
        println("client %d: %s", count, msg) // received from client <count>
        n, err := conn.Write([]byte("reply from server"))
        if err != nil {
            log.Println(n, err)
            return
        }
    }
}
