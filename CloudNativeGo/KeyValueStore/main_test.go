package main

// import (
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"strings"
// 	"testing"
// 	"time"
// )

// func TestAcceptance(t *testing.T) {
// 	var key = "key" + time.Now().String()
// 	const value = "value"
// 	const filename = "log.txt"

// 	s := httptest.NewServer(GetRouter())
// 	defer func() {
// 		s.Close()
// 		if err := os.Remove(filename); err != nil {
// 			t.Fatal("Could not cleanup file")
// 		}
// 	}()

// 	resp, err := s.Client().Get(s.URL + "/v1/key/" + key)
// 	if err != nil {
// 		t.Fatal("Error making GET request")
// 	}
// 	if resp.StatusCode != http.StatusNotFound {
// 		t.Fatal("Key found before created")
// 	}

// 	req, err := http.NewRequest("PUT", s.URL+"/v1/key/"+key, strings.NewReader(value))
// 	if err != nil {
// 		t.Fatal("Can't create PUT request")
// 	}
// 	resp, err = s.Client().Do(req)
// 	if err != nil || resp.StatusCode != http.StatusCreated {
// 		t.Fatal("Err making PUT request")
// 	}

// 	resp, err = s.Client().Get(s.URL + "/v1/key/" + key)
// 	if err != nil {
// 		t.Fatal("Error making GET request")
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		t.Fatal("Key not found after create")
// 	}
// 	respValue, err := io.ReadAll(resp.Body)
// 	defer resp.Body.Close()
// 	if err != nil {
// 		t.Fatal("Cannot read body")
// 	}
// 	if string(respValue) != value {
// 		t.Fatal("Value does not match")
// 	}

// 	req, err = http.NewRequest("DELETE", s.URL+"/v1/key/"+key, nil)
// 	if err != nil {
// 		t.Fatal("Can't create DELETE request")
// 	}
// 	resp, err = s.Client().Do(req)
// 	if err != nil || resp.StatusCode != http.StatusOK {
// 		t.Fatal("Err making DELETE request")
// 	}

// 	resp, err = s.Client().Get(s.URL + "/v1/key/" + key)
// 	if err != nil {
// 		t.Fatal("Error making GET request")
// 	}
// 	if resp.StatusCode != http.StatusNotFound {
// 		t.Fatal("Key found after delete")
// 	}
// }

// // TODO
// // this doesn't work with reading a log file that has no value for delete (even though value isn't used)
// func TestCanReadFile(t *testing.T) {
// 	const filename = "log.txt"
// 	const testfilename = "test_log.txt"
// 	var keys = []string{"key1", "key2", "key3"}

// 	file, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
// 	if err != nil {
// 		t.Fatal("Could not open log file")
// 	}
// 	defer file.Close()

// 	testfile, err := os.OpenFile(testfilename, os.O_RDONLY, 0755)
// 	if err != nil {
// 		t.Fatal("Could not open test file")
// 	}
// 	defer testfile.Close()

// 	_, err = io.Copy(file, testfile)
// 	if err != nil {
// 		t.Fatal("Could not setup log transactions")
// 	}

// 	s := httptest.NewServer(GetRouter())
// 	defer func() {
// 		s.Close()
// 		if err := os.Remove(filename); err != nil {
// 			t.Fatal("Could not cleanup file")
// 		}
// 	}()
// 	validateKey(t, s, false, keys[0], "")
// 	validateKey(t, s, true, keys[1], keys[1])
// 	validateKey(t, s, true, keys[2], keys[2])
// }

// func validateKey(t *testing.T, s *httptest.Server, exists bool, key string, value string) {
// 	if exists {
// 		resp, err := s.Client().Get(s.URL + "/v1/key/" + key)
// 		if err != nil {
// 			t.Fatal("Error making GET request")
// 		}
// 		if resp.StatusCode != http.StatusOK {
// 			t.Fatal("Key not found after create")
// 		}
// 		respValue, err := io.ReadAll(resp.Body)
// 		defer resp.Body.Close()
// 		if err != nil {
// 			t.Fatal("Cannot read body")
// 		}
// 		if string(respValue) != value {
// 			t.Fatal("Value does not match")
// 		}
// 	} else {
// 		resp, err := s.Client().Get(s.URL + "/v1/key/" + key)
// 		if err != nil {
// 			t.Fatal("Error making GET request")
// 		}
// 		if resp.StatusCode != http.StatusNotFound {
// 			t.Fatal("Found key that should not exist")
// 		}
// 	}
// }
