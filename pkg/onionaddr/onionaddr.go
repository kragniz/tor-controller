package onionaddr

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base32"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
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

func GetAddress(key []byte) (string, error) {
	// load the private key
	private, err := LoadPrivateKey(key)
	if err != nil {
		return "", nil
	}

	// get the public key
	public := private.PublicKey
	publicBytes, err := x509.MarshalPKIXPublicKey(&public)
	if err != nil {
		return "", err
	}

	// strip the first 22 bytes
	trimmedBytes := publicBytes[22:]

	// take the SHA1 digest
	sha := sha1.New()
	_, err = sha.Write(trimmedBytes)
	if err != nil {
		return "", err
	}
	digest := sha.Sum(nil)

	// base32 encode it
	b32 := base32.StdEncoding.EncodeToString(digest)

	// lowercase and take the first 16 characters
	return fmt.Sprintf("%s.onion", strings.ToLower(b32[:16])), nil
}
