package onionaddr

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func LoadPrivateKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("couldn't load key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}
