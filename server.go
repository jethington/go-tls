package main

// todo not sure if all of these are needed?
import (
    //"crypto/rand"
    "crypto/tls"
    "log"
    "net"
    "bufio"
    "crypto/x509"
    "io/ioutil"
)

// References:
//  https://golang.org/pkg/crypto/tls/
//  https://gist.github.com/jim3ma/00523f865b8801390475c4e2049fe8c3
//  https://gist.github.com/denji/12b3a568f092ab951456
//  https://ericchiang.github.io/post/go-tls/

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
    //fmt.Print(string(dat))
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
        SessionTicketsDisabled:      true,
        ClientAuth:                 tls.RequireAndVerifyClientCert,
    }

    ln, err := tls.Listen("tcp", "127.0.0.1:8883", &cfg) // note, example assigns config := &, and then passes config
    if err != nil {
        //log.Println(err)
        //return
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
            return // TODO is this really needed?
        }
    }
}

