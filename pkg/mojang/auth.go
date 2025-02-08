package mojang

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

type KeyPair struct {
	private *rsa.PrivateKey
	public  []byte
}

func NewKeyPairFromPrivateKey(privateKey *rsa.PrivateKey) (KeyPair, error) {
	if err := privateKey.Validate(); err != nil {
		return KeyPair{}, err
	}

	publicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return KeyPair{}, err
	}

	return KeyPair{privateKey, publicKey}, nil
}

func GenerateKeyPair() (KeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return KeyPair{}, err
	}
	return NewKeyPairFromPrivateKey(privateKey)
}

func (kp KeyPair) PublicKey() []byte {
	return kp.public
}

func (kp KeyPair) Decrypt(buffer []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(nil, kp.private, buffer)
}
