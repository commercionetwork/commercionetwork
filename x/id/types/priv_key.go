package types

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func ParseRsaPrivKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, fmt.Errorf("key is not pem encoded")
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, fmt.Errorf("key is not rsa priv key")
	}

	return pkey, nil
}

func PublicKeyToPemString(pub *rsa.PublicKey) string {
	return string(
		pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: x509.MarshalPKCS1PublicKey(pub),
			},
		),
	)
}

func LoadRSAPrivKeyFromDisk(location string) (*rsa.PrivateKey, error) {
	keyData, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}
	key, err := ParseRsaPrivKeyFromPEM(keyData)
	if err != nil {
		return nil, err
	}

	return key, nil
}
