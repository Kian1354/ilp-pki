package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"time"
	"reflect"

	"github.com/go-redis/redis"
)

// Initializes a Redis client to store certificates.
var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func main() {
	pong, _ := client.Ping().Result()
	fmt.Println("Doing health check for redis client: ", pong)

	// issue an identity for ILP address "123"
	privKey := IssueIdentity("123")
	fmt.Println("Generated Private Key of type:", reflect.TypeOf(privKey))

	// access public key elements
	// privKey.PublicKey.N --> big modulus
	// privKey.PublicKey.E --> public exponent

	// access private key elements
	// privKey.D --> private exponent
	// privKey.Primes --> prime factors of N


	res := RetrieveIdentity("123")

	fmt.Println("Retrieved certificate, with public key", res.PublicKey)
}

// Certificate struct contains the public key, ILP address, and the expiration time.
type Certificate struct {
	PublicKey  rsa.PublicKey
	ILPAddress string
	Expiration time.Time
}

// Given a ILP address, generate and store a new certificate; return the private key.
func IssueIdentity(ILPAddress string) *rsa.PrivateKey {

	// generate a 2048 bit RSA key
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)

	// gets the public key portion as a *rsa.PublicKey object
	// and create a certificate
	pub := priv.PublicKey
	c := newCertificate(pub, ILPAddress)

	// marshal the certificate and store in the Redis client
	bytes, _ := json.Marshal(c)
	client.Set(ILPAddress, bytes, 0)

	// return the private key
	return priv
}

// Given an ILP address, retrieve the certificate from the client.
func RetrieveIdentity(ilpAddress string) Certificate {

	// get the marshaled bytes from the redis client
	res, _ := client.Get(ilpAddress).Result()

	// unmarshal the bytes into a certificate struct variable
	var c Certificate
	_ = json.Unmarshal([]byte(res), &c)

	return c
}

// Given a public key and an ILP address, generate a new certificate.
func newCertificate(pub rsa.PublicKey, ILPAddress string) Certificate {
	expiration := time.Now().AddDate(1, 0, 0)

	c := Certificate{pub, ILPAddress, expiration}

	return c
}

// TBD
func reissueIdentity() bool {
	return false
}

// TBD
func revokeIdentity() bool {
	return false
}
