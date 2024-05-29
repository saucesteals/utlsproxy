module github.com/saucesteals/utlsproxy

go 1.21.0

// replace github.com/elazarl/goproxy => github.com/saucesteals/goproxy v0.0.0-20240124022437-840670a451ca
replace github.com/elazarl/goproxy => ../goproxy

require (
	github.com/elazarl/goproxy v0.0.0-20240124022437-840670a451ca
	github.com/refraction-networking/utls v1.6.0
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/cloudflare/circl v1.3.6 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/quic-go/quic-go v0.37.4 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)
