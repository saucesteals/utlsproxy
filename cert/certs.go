package cert

import (
	"crypto/tls"
	"crypto/x509"
)

func GetCertificate() (*tls.Certificate, error) {
	if err := setup(); err != nil {
		return nil, err
	}

	ca, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	if ca.Leaf, err = x509.ParseCertificate(ca.Certificate[0]); err != nil {
		return nil, err
	}

	return &ca, nil
}
