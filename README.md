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
  -addr string
        Address to bind to (default ":8080")
  -keylog string
        TLS key log file
  -http1
        Force HTTP/1.1 between client and proxy
```

---

## Why?

All (to my knowledge) MITM proxies replay requests to servers with stdlib transports, essentially letting the server fingerprint it. The goal of utlsproxy is to allow you to inspect TLS application data while mimicking the original client. The proxy will sniff the client's clienthello message and use it as its own when handshaking with the server.

Curious how? Most of the work is at [saucesteals/goproxy](https://github.com/saucesteals/goproxy) (credits to [elazarl/goproxy](https://github.com/elazarl/goproxy) for the base proxy implementation)

## Injecting a Client Hello from a previous session

Instead of fingerprinting the proxy client's ClientHello, you might want to save a ClientHello and re-inject it. E.g. you can save a Safari ClientHello and use it for your cURL requests.

### Saving a Client Hello

Simply define the `GOPROXY_CLIENT_HELLO_SAVE_DIR` variable:

```bash
GOPROXY_CLIENT_HELLO_SAVE_DIR="./client_hello" ./utlsproxy
```

This will save the client hello files in the `./client_hello` directory.

### Re-using a saved Client Hello

This time, define the `GOPROXY_OVERWRITE_CLIENT_HELLO` variable:

```bash
GOPROXY_OVERWRITE_CLIENT_HELLO="./client_hello/ch_safari_17.4.1_macOS_14.4.1.bin" ./utlsproxy
```

All requests will then have Safari's fingerprint.

To confirm

```bash
curl --silent --insecure --proxy localhost:8080 https://tls.peet.ws/api/tls | jq .tls.peetprint_hash
# "b2bafdc69377086c3416be278fd21121"
```

## mTLS

Like every other MITM, this will not work with mTLS. Find the client's certificate and private key, then add it to the tls.Config (Rarely will you need this, so this is only possible by cloning and adding it yourself)

## Contributing

Contributions are welcome!

- **[Submit Pull Requests](https://github.com/saucesteals/utlsproxy/pulls)**
- **[Report Issues](https://github.com/saucesteals/utlsproxy/issues)**

## License

Distributed under the GNU GPL v3.0 License. See `LICENSE` for more information.
