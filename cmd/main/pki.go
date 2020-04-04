package main

import (
	"bufio"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/go-redis/redis"
)

var YEAR = 1
var MONTH = 0
var DAY = 0
var PATH = "./gen_ca/ca_credentials/ca.key"

// Initializes a Redis client to store certificates.
var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

// TestFunctionality is a dummy method that shows how to interact with the different methods in pki.go
func TestFunctionality() {
	pong, _ := client.Ping().Result()
	fmt.Println("Doing health check for redis client: ", pong)

	// issue an identity for ILP address "123"
	privKey, _ := IssueIdentity("567")
	fmt.Println("Generated Private Key of type:", reflect.TypeOf(privKey))

	// access public key elements
	// privKey.PublicKey.N --> big modulus
	// privKey.PublicKey.E --> public exponent

	// access private key elements
	// privKey.D --> private exponent
	// privKey.Primes --> prime factors of N

	// retrieve the corresponding identity for ILP address "123"
	res, _ := RetrieveIdentity("567")
	fmt.Println("Retrieved certificate, with public key", res.Cert.PublicKey)

	// save the SignedCertificate to a file
	SaveCertificateToFile(res)

	// verify that the SignedCertificate is valid
	valid := IsValidCertificate(res)
	fmt.Println("Certificate is valid: ", valid)
}

// SignedCertificate stores a certificate and its signature.
type SignedCertificate struct {
	Cert      Certificate
	Signature []byte
}

// Certificate struct contains the public key, ILP address, and the expiration time.
type Certificate struct {
	PublicKey  rsa.PublicKey
	ILPAddress string
	Expiration time.Time
}

// IssueIdentity generates and stores a new certificate given a new ILP address; returns the private key. Given a ILP address, generate and store a new certificate; return the private key.
func IssueIdentity(ILPAddress string) (*rsa.PrivateKey, error) {
	// generate a 2048 bit RSA key
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)

	// gets the public key portion as a *rsa.PublicKey object
	// and create a certificate
	pub := priv.PublicKey
	c := newCertificate(pub, ILPAddress)

	signedC := newSignedCertificate(c)

	// marshal the signed certificate and store in the Redis client
	bytes, _ := json.Marshal(signedC)

	expiration := c.Expiration
	duration := time.Until(expiration)

	client.Set(ILPAddress, bytes, duration)

	// return the private key
	return priv, nil
}

// RetrieveIdentity retrieves a certificate for a given ILP address.
func RetrieveIdentity(ilpAddress string) (SignedCertificate, error) {

	cert := &SignedCertificate{}
	// get the marshaled bytes from the redis client
	res, err := client.Get(ilpAddress).Result()

	if err != nil {
		return *cert, err
	}

	// unmarshal the bytes into a certificate struct variable
	// var c Certificate
	_ = json.Unmarshal([]byte(res), cert)

	return *cert, nil
}

// Given an ILP address, remove the certificate from the DB.
func revokeIdentity(ilpAddress string) error {
	return client.Del(ilpAddress).Err()
}

// IsValidCertificate checks if a given certificate is valid is valid; certificate is defined to be valid if signature is valid and certificate is not expired.
func IsValidCertificate(signedC SignedCertificate) bool {
	return signedC.Cert.Expiration.After(time.Now()) && verifySignature(signedC.Cert, signedC.Signature)
}

// SaveCertificateToFile saves a given certificate to a file.
func SaveCertificateToFile(certificate SignedCertificate) (fileName string, err error) {
	jsonData, err := json.Marshal(certificate)

	// write to JSON file
	jsonFile, err := os.Create("./certificate_files/certificate" + (certificate.Cert.ILPAddress)[0:3] + ".json")

	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()
	fmt.Println("JSON data written to ", jsonFile.Name())

	return jsonFile.Name(), err

}

// Given a public key and an ILP address, generate a new certificate.
func newCertificate(pub rsa.PublicKey, ILPAddress string) Certificate {
	expiration := time.Now().AddDate(YEAR, MONTH, DAY)
	c := Certificate{pub, ILPAddress, expiration}

	return c
}

func newSignedCertificate(cert Certificate) SignedCertificate {
	bytes, _ := json.Marshal(cert)

	sig := createSignature(bytes)

	signedCert := SignedCertificate{cert, sig}

	return signedCert
}

func createSignature(message []byte) []byte {
	rng := rand.Reader

	hash := sha256.Sum256(message)
	MSK, _ := getCAKey()

	sig, err := rsa.SignPKCS1v15(rng, MSK, crypto.SHA256, hash[:])

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
		return nil
	}

	return sig
}

// Returns true if signature successfully verified
func verifySignature(cert Certificate, signature []byte) bool {
	bytes, _ := json.Marshal(cert)

	hash := sha256.Sum256(bytes)

	MSK, _ := getCAKey()

	err := rsa.VerifyPKCS1v15(&MSK.PublicKey, crypto.SHA256, hash[:], signature)

	return err == nil
}

// // Returns true if signature successfully verified
// func verifyJSONCertificate(certFile []byte, signature []byte) bool {
// 	certificate := SignedCertificate{}
// 	err := json.Unmarshal([]byte(certFile), &certificate)

// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}
// 	isValid := IsValidCertificate(certificate)

// 	return isValid
// }

// Server method for retrieving the CA private key from file.
func getCAKey() (*rsa.PrivateKey, error) {

	privKeyFile, _ := os.Open(PATH)

	fileInfo, _ := privKeyFile.Stat()
	var size = fileInfo.Size()

	keyBytes := make([]byte, size)

	buffer := bufio.NewReader(privKeyFile)
	_, err := buffer.Read(keyBytes)

	data, _ := pem.Decode([]byte(keyBytes))
	privKeyFile.Close()

	MSK, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return MSK, nil
}
