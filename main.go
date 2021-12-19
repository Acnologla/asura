package main

import (
	"asura/src/server"
	"crypto/tls"
	"fmt"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

var Port string

func main() {
	if os.Getenv("PRODUCTION") == "" {
		err := godotenv.Load()
		if err != nil {
			panic("Cannot read the motherfucking envfile")
		}
	}
	server.Init(os.Getenv("PUBLIC_KEY"))
	Port = os.Getenv("PORT")

	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("localhost"), // Replace with your domain.
		Cache:      autocert.DirCache("./certs"),
	}

	cfg := &tls.Config{
		GetCertificate: m.GetCertificate,
		NextProtos: []string{
			"http/1.1", acme.ALPNProto,
		},
	}

	ln, err := net.Listen("tcp4", "0.0.0.0:"+Port) /* #nosec G102 */
	if err != nil {
		panic(err)
	}

	lnTls := tls.NewListener(ln, cfg)

	if err := fasthttp.Serve(lnTls, server.Handler); err != nil {
		panic(err)
	}

	fmt.Println("server started")
}
