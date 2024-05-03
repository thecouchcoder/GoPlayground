package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAcceptance(t *testing.T) {
	const key = "key"
	const value = "value"

	s := httptest.NewServer(GetRouter())
	defer s.Close()

	resp, err := s.Client().Get(s.URL + "/v1/key/" + key)
	if err != nil {
		t.Fatal("Error making GET request")
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatal("Key found before created")
	}

	req, err := http.NewRequest("PUT", s.URL+"/v1/key/"+key, strings.NewReader(value))
	if err != nil {
		t.Fatal("Can't create PUT request")
	}
	resp, err = s.Client().Do(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatal("Err making PUT request")
	}

	resp, err = s.Client().Get(s.URL + "/v1/key/" + key)
	if err != nil {
		t.Fatal("Error making GET request")
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Key not found after create")
	}
	respValue, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Fatal("Cannot read body")
	}
	if string(respValue) != value {
		t.Fatal("Value does not match")
	}

	req, err = http.NewRequest("DELETE", s.URL+"/v1/key/"+key, nil)
	if err != nil {
		t.Fatal("Can't create DELETE request")
	}
	resp, err = s.Client().Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatal("Err making DELETE request")
	}

	resp, err = s.Client().Get(s.URL + "/v1/key/" + key)
	if err != nil {
		t.Fatal("Error making GET request")
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatal("Key found after delete")
	}
}
