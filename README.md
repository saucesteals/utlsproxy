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
$ go install github.com/saucesteals/utlsproxy
```

---

## Usage

```sh
$ utlsproxy
  -interface string
        Network interface to bind to (default "lo0")
  -port string
        Proxy port (default "8080")
  -keylog string
        TLS key log file
```

---

## Why?

All (to my knowledge) MITM proxies replay requests to servers with stdlib transports, essentially letting the server fingerprint it. The goal of utlsproxy is to allow you to inspect TLS application data while mimicing the original client. The proxy will sniff the client's clienthello message and use it as its own when handshaking with the server.

Curious how? Most of the work is at [saucesteals/goproxy](https://github.com/saucesteals/goproxy) (credits to [elazarl/goproxy](https://github.com/elazarl/goproxy) for the base proxy implementation)

## mTLS

Like every other MITM, this will not work with mTLS. Find the client's certificate and private key, then add it to the tls.Config (Rarely will you need this, so this is only possible by cloning and adding it yourself)

## Contributing

Contributions are welcome!

- **[Submit Pull Requests](https://github.com/saucesteals/utlsproxy/pulls)**
- **[Report Issues](https://github.com/saucesteals/utlsproxy/issues)**

## License

Distributed under the GNU GPL v3.0 License. See `LICENSE` for more information.
