package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func main() {
	pong, err := client.Ping().Result()
	fmt.Print("Doing health check for redis client: ")
	fmt.Println(pong, err)

	err = client.Set("name", "Kian", 0).Err()
	res, err := client.Get("name").Result()

	fmt.Println(res)
	fmt.Println(time.Now().AddDate(1, 0, 0))
	// time.Now() = 2019-10-20 22:59:55.644117 -0700 PDT m=+1.734839423

	issueIdentity("vvl")
	res1, _ := client.Get("vvl").Result()
	fmt.Println(res1)
	var c Certificate
	_ = json.Unmarshal([]byte(res1), &c)
	fmt.Println(c.ILPAddress)
}

type Certificate struct {
	PublicKey  crypto.PublicKey
	ILPAddress string
	Expiration time.Time
}

func newCertificate(pub crypto.PublicKey, ILPAddress string) Certificate {
	expiration := time.Now().AddDate(1, 0, 0)

	c := Certificate{pub, ILPAddress, expiration}

	return c
}

func issueIdentity(ILPAddress string) rsa.PrivateKey {

	// generate a 2048 bit RSA key
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)

	// gets the public key portion as a crypto.PublicKey object
	// and create a certificate
	pub := priv.Public()
	c := newCertificate(pub, ILPAddress)

	// marshal the certificate and store
	bytes, _ := json.Marshal(c)
	client.Set(ILPAddress, bytes, 0)

	return *priv
}

func reissueIdentity() bool {
	return false
}

func revokeIdentity() bool {
	return false
}

// should return a certificate
func retrieveIdentity(ilpAddress string) Certificate {
	return false
}
