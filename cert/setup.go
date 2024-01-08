package cert

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

var (
	certFileName = "cert.pem"
	keyFileName  = "key.pem"

	cert, key []byte
)

func setup() error {
	if cert != nil && key != nil {
		return nil
	}

	if err := createKeyPair(); err != nil {
		return err
	}

	var err error
	cert, key, err = readKeyPair()
	if err != nil {
		return err
	}

	return nil
}

func prefsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	fp := filepath.Join(home, ".config", "utlsproxy")
	if err := os.MkdirAll(fp, 0700); err != nil {
		return "", err
	}

	return fp, nil
}

func createKeyPair() error {
	fp, err := prefsDir()
	if err != nil {
		panic(err)
	}

	certFile := filepath.Join(fp, certFileName)
	keyFile := filepath.Join(fp, keyFileName)

	if _, err := os.Stat(certFile); err == nil {
		return nil
	}

	user, err := user.Current()
	if err != nil {
		return err
	}

	subj := fmt.Sprintf("/O=utlsproxy/OU=%s/CN=utlsproxy", user.Username)
	if err := exec.Command("openssl", "req", "-x509", "-newkey", "rsa:2048", "-keyout", keyFile, "-out", certFile, "-sha256", "-days", "3650", "-nodes", "-subj", subj).Run(); err != nil {
		return err
	}

	return nil
}

func readKeyPair() ([]byte, []byte, error) {
	fp, err := prefsDir()
	if err != nil {
		panic(err)
	}

	certFile := filepath.Join(fp, certFileName)
	keyFile := filepath.Join(fp, keyFileName)

	cert, err := os.ReadFile(certFile)
	if err != nil {
		return nil, nil, err
	}

	key, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}

	return cert, key, nil
}
