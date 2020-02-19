package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

// Testing Suite
func TestIsValidCertificate(t *testing.T) {
	_, err := IssueIdentity("123456789")
	if err != nil {
		t.Errorf("Error issuing identity")
	}

	certificate, err := RetrieveIdentity("123456789")
	if err != nil {
		t.Errorf("Error retrieving identity")
	}

	filepath, err := SaveCertificateToFile(certificate)
	if err != nil {
		t.Errorf("Error saving certificate to file")
	}

	certificateFile, _ := ioutil.ReadFile(filepath)

	certificate = SignedCertificate{}
	_ = json.Unmarshal([]byte(certificateFile), &certificate)

	isValid := IsValidCertificate(certificate)
	if !isValid {
		t.Errorf("Certificate should be valid but is not")
	}

}

func TestIsValidCertificateFail(t *testing.T) {
	_, err := IssueIdentity("123456789")
	if err != nil {
		t.Errorf("Error issuing identity")
	}

	certificate, err := RetrieveIdentity("123456789")
	if err != nil {
		t.Errorf("Error retrieving identity")
	}

	filepath, err := SaveCertificateToFile(certificate)
	if err != nil {
		t.Errorf("Error saving certificate to file")
	}

	certificateFile, _ := ioutil.ReadFile(filepath)

	certificate = SignedCertificate{}
	_ = json.Unmarshal([]byte(certificateFile), &certificate)

	certificate.Signature = []byte("876543457898765434567876543")

	isValid := IsValidCertificate(certificate)
	if isValid {
		t.Errorf("Certificate should not be valid but it is")
	}

}

func TestRevokeIdentity(t *testing.T) {
	_, err := IssueIdentity("123456789")
	if err != nil {
		t.Errorf("Error issuing identity")
	}

	_, err = RetrieveIdentity("123456789")
	if err != nil {
		t.Errorf("Error retrieving identity")
	}

	err = revokeIdentity("123456789")

	if err != nil {
		t.Errorf("Error revoking identity")
	}

	_, err = RetrieveIdentity("123456789")
	if err == nil {
		t.Errorf("Incorrectly retrieved identity after revoking")
	}

}
