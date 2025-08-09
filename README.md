<div align="center">
<img src="./assets/logo.png" alt="logo" height="250px" />
<h1 align="center">uTLS Proxy</h1>
<h3>MITM Proxy with TLS mimicry</h3>
</div>

---

## Features

- Man-in-the-middle TLS connections
- Mimic TLS records from the client
- Directly proxy all application data below TLS (for example, http2 frames will be exactly what the client sends)
- SSL/TLS Keylog support (inspect TLS contents in wireshark)
- Certificate server at `utlsproxy.ws` (certificates are created locally)

---

## Installation

```sh
$ git clone https://github.com/saucesteals/utlsproxy
$ cd utlsproxy
$ go build
$ ./utlsproxy
```

---

## Usage

```sh
$ ./utlsproxy
  -addr string
        Address to bind to (default ":8080")
  -http1
        Force HTTP/1.1 between client and proxy
  -keylog string
        TLS key log file
  -clientcert string
        mTLS client certificate file (pem)
  -clientkey string
        mTLS client key file (pem)
  -mtlsdomain string
        Enable mTLS for this domain
```

---

## Why?

All (to my knowledge) MITM proxies replay requests to servers with stdlib transports, essentially letting the server fingerprint it. The goal of utlsproxy is to allow you to inspect TLS application data while mimicking the original client. The proxy will sniff the client's clienthello message and use it as its own when handshaking with the server.

Curious how? Most of the work is at [saucesteals/goproxy](https://github.com/saucesteals/goproxy) (credits to [elazarl/goproxy](https://github.com/elazarl/goproxy) for the base proxy implementation)

## Contributing

Contributions are welcome!

- **[Submit Pull Requests](https://github.com/saucesteals/utlsproxy/pulls)**
- **[Report Issues](https://github.com/saucesteals/utlsproxy/issues)**

## License

Distributed under the GNU GPL v3.0 License. See `LICENSE` for more information.
