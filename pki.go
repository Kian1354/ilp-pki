package main

import (
	"fmt"
	"crypto/rand"
	"crypto/rsa"
	"time"
	"encoding/json"
	"github.com/go-redis/redis"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "password",
	DB:       0,
})

func main() {
	fmt.Println("hi there!")

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	err = client.Set("name", "Kian", 0).Err()
	res, err := client.Get("name").Result()

	fmt.Println(res)
	fmt.Println(time.Now().AddDate(1,0,0))
	// time.Now() = 2019-10-20 22:59:55.644117 -0700 PDT m=+1.734839423

	issueIdentity("vvl")
	res1, _ := client.Get("vvl").Result()
	fmt.Println(res1)
	// var c Certificate
	// unmarshalled := json.Unmarshal([]byte(res1), c)
	// fmt.Println(unmarshalled)
}


type Certificate struct {
	PublicKey rsa.PublicKey
	ILPAddress string
	Expiration time.Time
}

func newCertificate(pub rsa.PublicKey, ILPAddress string) Certificate {
	expiration := time.Now().AddDate(1,0,0)

	c := Certificate{pub, ILPAddress, expiration}

	// bytes, err := json.Marshal(c)
	// fmt.Println(bytes)
	// client.set(ILPAddress, )
	return c
}

func issueIdentity(ILPAddress string) bool {

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey

	// fmt.Println(*pub)
	c := newCertificate(*pub, ILPAddress)

	bytes, _ := json.Marshal(c)
	// fmt.Println(bytes)

	client.Set(ILPAddress, bytes, 0)

	return true
}

func reissueIdentity() bool {
	return false
}

func revokeIdentity() bool {
	return false
}

// should return a certificate
func retrieveIdentity(ilpAddress string) bool {
	return false
}
