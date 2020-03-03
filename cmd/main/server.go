package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type createIdentityRequest struct {
	IlpAddress string
}

type createIdentityResponse struct {
	PrivKey []byte `json:"privKey"`
}

type retrieveIdentityRequest struct {
	IlpAddress string
}

type retrieveIdentityResponse struct {
	SignedCert []byte `json:"signedCert"`
}

func createIdentity(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		fmt.Println("POST request required.")
	}

	decoder := json.NewDecoder(r.Body)
	var request createIdentityRequest
	err := decoder.Decode(&request)
	if err != nil {
		fmt.Println("Decoding request failed.")
	}

	privKey, err := IssueIdentity(request.IlpAddress)
	if err != nil {
		fmt.Println("Issue Identity failed.")
	}

	bytes, _ := json.Marshal(privKey)
	w.Header().Set("Content-Type", "application/json")
	response := createIdentityResponse{
		PrivKey: bytes,
	}
	json.NewEncoder(w).Encode(response)
}

func retrieveIdentity(w http.ResponseWriter, r *http.Request) {
	// returns a json object corresponding to a certificate

	if r.Method != "GET" {
		fmt.Println("GET request required.")
	}

	decoder := json.NewDecoder(r.Body)
	var request retrieveIdentityRequest
	err := decoder.Decode(&request)
	if err != nil {
		fmt.Println("Decoding request failed.")
	}

	signedCert, err := RetrieveIdentity(request.IlpAddress)
	if err != nil {
		fmt.Println("Retrieve Identity failed.")
	}

	bytes, _ := json.Marshal(signedCert)
	w.Header().Set("Content-Type", "application/json")
	response := retrieveIdentityResponse{
		SignedCert: bytes,
	}
	json.NewEncoder(w).Encode(response)
}

func verifyIdentity(w http.ResponseWriter, r *http.Request) {
	// a marshalled representation of a json file representing the certificate is passed into r.Body
	decoder := json.NewDecoder(r.Body)
	var cert SignedCertificate
	err := decoder.Decode(&cert)

	if err != nil {
		fmt.Fprintf(w, "Could not decode")
	}

	isValid := IsValidCertificate(cert)
	fmt.Fprintf(w, strconv.FormatBool(isValid))
}

func main() {

	c, _ := RetrieveIdentity("kian")
	SaveCertificateToFile(c)

	// Create Identity
	http.HandleFunc("/create", createIdentity)

	// Retrieve Identity
	http.HandleFunc("/retrieve", retrieveIdentity)

	// Verify Certificate
	http.HandleFunc("/verify", verifyIdentity)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
