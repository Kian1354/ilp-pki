package main

import (
	"bufio"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const trustedCertURL = "https://raw.githubusercontent.com/Kian1354/ilp-pki/master/cmd/gen_ca/ca_credentials/ca.crt"

const trustedCertKeyURL = "https://raw.githubusercontent.com/Kian1354/ilp-pki/master/cmd/gen_ca/ca_credentials/ca.key"

const PATH = "./root_ca_cert/ca.crt"

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

type createIdentityResponse struct {
	PrivKey []byte `json:"privKey"`
}

type retrieveIdentityResponse struct {
	SignedCert []byte `json:"signedCert"`
}

// creates an identity associated with ilpAddress
func createIdentity(ilpAddress string, url string) *rsa.PrivateKey {

	url = url + "/create"

	s := fmt.Sprintf("{\n\t\"ilpAddress\": \"%s\"\n}", ilpAddress)
	payload := strings.NewReader(s)

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	responseObject := &createIdentityResponse{}
	err := json.Unmarshal(body, responseObject)

	if err != nil {
		fmt.Printf("Response error")
	}

	privKey := &rsa.PrivateKey{}
	_ = json.Unmarshal(responseObject.PrivKey, &privKey)

	return privKey
}

// retrieves the signed certificate associated with ilpAddress
func retrieveIdentity(ilpAddress string, url string) SignedCertificate {

	url = url + "/retrieve"

	s := fmt.Sprintf("{\n\t\"ilpAddress\": \"%s\"\n}", ilpAddress)
	payload := strings.NewReader(s)

	req, _ := http.NewRequest("GET", url, payload)
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	responseObject := &retrieveIdentityResponse{}
	err := json.Unmarshal(body, responseObject)

	if err != nil {
		fmt.Printf("Response error")
	}

	certificate := SignedCertificate{}
	_ = json.Unmarshal(responseObject.SignedCert, &certificate)

	return certificate
}

func getCAKey() (*rsa.PublicKey, error) {

	pubKeyFile, err := os.Open(PATH)
	if err != nil {
		fmt.Println("Please bootstrap CA key")
		return &rsa.PublicKey{}, err
	}

	fileInfo, _ := pubKeyFile.Stat()
	var size = fileInfo.Size()

	keyBytes := make([]byte, size)

	buffer := bufio.NewReader(pubKeyFile)
	_, err = buffer.Read(keyBytes)

	data, _ := pem.Decode([]byte(keyBytes))

	pubKeyFile.Close()

	MPK, err := x509.ParseCertificate(data.Bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return MPK.PublicKey.(*rsa.PublicKey), nil
}

func verifySignature(signedCert SignedCertificate) bool {
	cert := signedCert.Cert
	signature := signedCert.Signature

	bytes, _ := json.Marshal(cert)

	hash := sha256.Sum256(bytes)

	MPK, err := getCAKey()
	if err != nil {
		return false
	}

	err = rsa.VerifyPKCS1v15(MPK, crypto.SHA256, hash[:], signature)

	return err == nil
}

// DownloadFile will download a url and store it in local filepath.
// It writes to the destination file as it downloads it, without
// loading the entire file into memory.
//
// Source: https://gist.github.com/cnu/026744b1e86c6d9e22313d06cba4c2e9
func downloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// Helper Method to export certificate to file
// SaveCertificateToFile saves a given certificate to a file.
func SaveCertificateToFile(certificate SignedCertificate) (fileName string, err error) {
	jsonData, err := json.Marshal(certificate)

	// write to JSON file
	path := "./certificate_files"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0700)
	}
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

// Bootstrap gets trusted root CA cert from Github and puts cert in Folder
func Bootstrap() {
	path := "./root_ca_cert"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0700)
	}
	downloadFile(trustedCertURL, "./root_ca_cert/ca.crt")
	downloadFile(trustedCertKeyURL, "./root_ca_cert/ca.key")
}

func main() {
	Bootstrap()
	createIdentity("kian", "http://localhost:8080")
	res := retrieveIdentity("kian", "http://localhost:8080")

	fmt.Println(verifySignature(res))
}
