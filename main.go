package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/elazarl/goproxy"
	utls "github.com/refraction-networking/utls"
	"github.com/saucesteals/utlsproxy/cert"
)

var (
	flagKeyLogFile = flag.String("keylog", "", "TLS key log file")
	flagPort       = flag.String("port", "8080", "Proxy port")
	flagInterface  = flag.String("interface", "lo0", "Network interface to bind to")
)

func main() {
	flag.Parse()

	ca, err := cert.GetCertificate()
	if err != nil {
		log.Panic(err)
	}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(ca)}

	proxy := goproxy.NewProxyHttpServer(tlsConfig())
	proxy.Verbose = true
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(serveCertificate(ca))

	addr, err := bindToAddr()
	if err != nil {
		log.Panic(err)
	}

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
		defer w.Close()

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

func bindToAddr() (string, error) {
	addr := ""
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Name != *flagInterface {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, ifaceAddr := range addrs {
			if ip, ok := ifaceAddr.(*net.IPNet); ok {
				addr = ip.IP.String()
				break
			}
		}
	}

	if addr == "" {
		return "", fmt.Errorf("could not find interface %s", *flagInterface)
	}

	return fmt.Sprintf("%s:%s", addr, *flagPort), nil
}
