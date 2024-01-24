package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/elazarl/goproxy" // github.com/saucesteals/goproxy
	utls "github.com/refraction-networking/utls"
	"github.com/saucesteals/utlsproxy/cert"
)

var (
	flagKeyLogFile = flag.String("keylog", "", "TLS key log file")
	flagAddr       = flag.String("addr", ":8080", "Address to bind to")
	flagHttp1Only  = flag.Bool("http1", false, "Force HTTP/1.1 between client and proxy")
)

func main() {
	flag.Parse()
	flag.Usage()

	ca, err := cert.GetCertificate()
	if err != nil {
		log.Panic(err)
	}

	proxyCAConfig := goproxy.TLSConfigFromCA(ca)
	proxyTlsConfig := func(host string, ctx *goproxy.ProxyCtx) (*tls.Config, error) {
		config, err := proxyCAConfig(host, ctx)
		if err != nil {
			return nil, err
		}

		config.NextProtos = []string{"h2", "http/1.1"}
		if *flagHttp1Only {
			config.NextProtos = []string{"http/1.1"}
		}

		return config, nil
	}

	proxy := goproxy.NewProxyHttpServer(tlsConfig())
	proxy.OnRequest().HandleConnect(goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		return &goproxy.ConnectAction{
			Action:    goproxy.ConnectMitm,
			TLSConfig: proxyTlsConfig,
		}, host
	}))
	proxy.OnRequest().DoFunc(serveCertificate(ca))

	addr := *flagAddr
	log.Printf("Listening on %s", addr)
	if err := http.ListenAndServe(addr, proxy); err != nil {
		log.Panic(err)
	}
}

func tlsConfig() *utls.Config {
	var keyLogWriter io.Writer
	if path := *flagKeyLogFile; path != "" {
		w, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			log.Panic(err)
		}

		keyLogWriter = w
	}

	return &utls.Config{KeyLogWriter: keyLogWriter}
}

func serveCertificate(ca *tls.Certificate) func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		certServer := "utlsproxy.ws"
		if r.URL.Host != certServer {
			return r, &http.Response{
				StatusCode: http.StatusNotFound,
				Body: io.NopCloser(
					strings.NewReader(fmt.Sprintf("<html><body>Not found. Did you mean <a href='http://%s'>%s</a> ?</body></html>", certServer, certServer)),
				),
			}
		}

		return r, &http.Response{
			StatusCode: http.StatusOK,
			Header: http.Header{
				"content-type": {"application/x-x509-ca-cert"},
			},
			Body: io.NopCloser(bytes.NewReader(ca.Leaf.Raw)),
		}
	}
}
