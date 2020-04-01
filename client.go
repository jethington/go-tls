package main

// todo not sure if all of these are needed?
import (
    "crypto/tls"
    "log"
    "crypto/x509"
    "io/ioutil"
    //"fmt"
)

// References:
//  https://golang.org/pkg/crypto/tls/
//  https://gist.github.com/jim3ma/00523f865b8801390475c4e2049fe8c3
//  https://gist.github.com/denji/12b3a568f092ab951456
//  https://ericchiang.github.io/post/go-tls/

func main() {
    cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
    if err != nil {
        log.Println(err)
        return
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

    config := tls.Config{
        ServerName:         "server",
        Certificates:       []tls.Certificate{cert},
        RootCAs:            roots,
        MinVersion:         tls.VersionTLS12,
        MaxVersion:         tls.VersionTLS12,
    }

    conn, err := tls.Dial("tcp", "127.0.0.1:8883", &config)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()

    n, err := conn.Write([]byte("client sending\n"),)
    if err != nil {
        log.Println(n, err)
        return
    }

    buf := make([]byte, 100) // TODO is this really the way to declare an array?
    n, err = conn.Read(buf)
    if err != nil {
        log.Println(n, err)
        return
    }

    println(string(buf[:n])) // TODO that's how you turn an array into a string?
}
