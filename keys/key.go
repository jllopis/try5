package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"time"
)

type Key struct {
	ID        *int64     `json:"id" db:"id"`
	KID       *string    `json:"kid" db:"kid"`
	PubKey    []byte     `json:"pubkey" db:"pub_key"`
	PrivKey   []byte     `json:"privkey" db:"priv_key"`
	AccountID *string    `json:"account_id" db:"account_id"`
	Active    *bool      `json:"active" db:"active"`
	Created   *time.Time `json:"created" db:"created"`
	Updated   *time.Time `json:"updated" db:"updated"`
	Deleted   *time.Time `json:"deleted,omitempty" db:"deleted"`
}

// New genera una pareja de claves RSA de 2048 bits. Las clave privada se codifica como PKCS1 y la pública como PKIX.
// Ambas en formato PEM.
func New(uid string) *Key {
	if uid == "" {
		log.Printf("new key needs an account")
		return nil
	}
	k := newKey()
	k.AccountID = &uid
	return k
}

// newKey realiza la generación y codificación de las claves RSA en codificación PEM.
func newKey() *Key {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Printf("failed to generate private key: %s", err)
		return nil
	}
	privPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
	pubKeyPKIX, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Printf("failed to generate DER public key: %s", err)
		return nil
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyPKIX,
	})

	return &Key{
		PubKey:  pubPEM,
		PrivKey: privPEM,
	}
}
