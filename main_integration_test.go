package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	initializeBlockChain()
}

func TestHandleMessage_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/message", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handleMessage(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestHandleValidation_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/validate", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handleMessage(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestHandleConflictResolution_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/resolve", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handleMessage(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestHandleMessage_MissingFields(t *testing.T) {
	req, err := http.NewRequest("POST", "/message", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handleMessage(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	expectedMessage := "All fields are required"
	if !strings.Contains(rr.Body.String(), expectedMessage) {
		t.Errorf("message was not found in the response body: got %v", rr.Body.String())
	}
}

func TestHandleMessage_ValidTransaction(t *testing.T) {
	body := strings.NewReader("sender=Alice&receiver=Bob&message=Hello!")
	req, err := http.NewRequest("POST", "/message", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handleMessage(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	chainCount := len(blockchain.Chain)
	if chainCount != 2 {
		t.Errorf("handler failed to add transaction to the blockchain: got %v chains, want %v", chainCount, 2)
	}
}

func TestHandleValidation_ValidBlockchain(t *testing.T) {
	req, err := http.NewRequest("POST", "/validate", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handleValidation(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	expectedMessage := "Blockchain is valid."
	if !strings.Contains(rr.Body.String(), expectedMessage) {
		t.Errorf("message was not found in the response body: got %v", rr.Body.String())
	}
}

func TestHandleValidation_InValidBlockchain(t *testing.T) {
	req, err := http.NewRequest("POST", "/validate", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	blockchain.Chain[1].Hash = "invalid hash"

	handleValidation(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", rr.Code, http.StatusOK)
	}

	expected := "Blockchain is not valid."
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v, want %v", rr.Body.String(), expected)
	}
}

func TestMessengerHandler(t *testing.T) {
	w := httptest.NewRecorder()

	r := httptest.NewRequest("GET", "/", nil)

	messengerHandler(w, r)

	resp := w.Result()

	expectedStatusCode := http.StatusOK
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("messengerHandler failed: got status code %v, want %v", resp.StatusCode, expectedStatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("messengerHandler failed to read response body: %v", err)
	}

	expectedContent := "Blockchain Messenger"
	if !strings.Contains(string(body), expectedContent) {
		t.Errorf("messengerHandler failed: response body doesn't contain expected content '%s'", expectedContent)
	}
}

func TestHandleConflictResolution(t *testing.T) {
	req, err := http.NewRequest("POST", "/resolve", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	HandleConflictResolution(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	expectedMessage := "Conflict resolution completed."
	if !strings.Contains(rr.Body.String(), expectedMessage) {
		t.Errorf("message was not found in the response body: got %v", rr.Body.String())
	}
}
