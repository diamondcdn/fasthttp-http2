package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/dgrr/http2"
	"github.com/valyala/fasthttp"
)

func main() {
	debug := flag.Bool("debug", true, "Debug mode")
	flag.Parse()

	cert, priv, err := GenerateTestCertificate("localhost:8443")
	if err != nil {
		log.Fatalln(err)
	}

	s := &fasthttp.Server{
		ReadTimeout: time.Second * 3,
		Handler:     requestHandler,
		Name:        "http2 test",
	}
	err = s.AppendCertEmbed(cert, priv)
	if err != nil {
		log.Fatalln(err)
	}

	http2.ConfigureServer(s, http2.ServerConfig{
		Debug: *debug,
	})

	err = s.ListenAndServeTLS(":8443", "", "")
	if err != nil {
		log.Fatalln(err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	if ctx.Request.Header.IsPost() {
		fmt.Fprintf(ctx, "%s\n", ctx.Request.Body())
		return
	}

	if ctx.FormValue("long") == nil {
		fmt.Fprintf(ctx, "Hello 21th century!\n")
	} else {
		for i := 0; i < 1<<16; i++ {
			ctx.Response.AppendBodyString("A")
		}
	}
}
