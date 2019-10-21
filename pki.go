package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func main() {
	fmt.Println("hi there!")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	err = client.Set("name", "Kian", 0).Err()
	res, err := client.Get("name").Result()

	fmt.Println(res)
}

func issueIdentity() bool {
	return false
}

func reissueIdentity() bool {
	return false
}

func revokeIdentity() bool {
	return false
}

// should return a certificate
func retrieveIdentity(ilpAddress string) {

}
