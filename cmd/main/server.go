package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type CreateIdentityRequest struct {
	IlpAddress string
}

type createIdentityResponse struct {
	PrivKey []byte `json:"privKey"`
}

func createIdentity(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		fmt.Println("POST request required.")
	}

	fmt.Println("body", r.Body)

	decoder := json.NewDecoder(r.Body)

	var request CreateIdentityRequest

	err := decoder.Decode(&request)
	fmt.Println("fuck you")
	fmt.Println(request.IlpAddress)

	if err != nil {
		panic(err)
	}

	privKey, err := IssueIdentity(request.IlpAddress)

	if err != nil {
		panic(err)
	}

	bytes, _ := json.Marshal(privKey)

	w.Header().Set("Content-Type", "application/json")
	response := createIdentityResponse{
		PrivKey: bytes,
	}

	// w.Write(bytes)
	json.NewEncoder(w).Encode(response)
}

func retrieveIdentity(w http.ResponseWriter, r *http.Request) {
	// returns a json object corresponding to a certificate

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
