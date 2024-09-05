package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCreatePost(t *testing.T) {
	// Test for x-www-form-urlencoded request
	t.Run("x-www-form-urlencoded", func(t *testing.T) {
		form := url.Values{}
		form.Add("h", "entry")
		form.Add("content", "This is a test post")
		form.Add("category[]", "test")
		form.Add("category[]", "example")

		req, err := http.NewRequest("POST", "/micropub", strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handleMicropub)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		expected := "Post created successfully"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	// Test for JSON request
	t.Run("JSON", func(t *testing.T) {
		jsonBody := map[string]interface{}{
			"type":    []string{"h-entry"},
			"content": "This is a test post",
			"category": []string{
				"test",
				"example",
			},
		}
		jsonBytes, _ := json.Marshal(jsonBody)

		req, err := http.NewRequest("POST", "/micropub", bytes.NewBuffer(jsonBytes))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handleMicropub)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		expected := "Post created successfully"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	// Test for unsupported content type
	t.Run("Unsupported Content-Type", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/micropub", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "text/plain")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handleMicropub)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnsupportedMediaType {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnsupportedMediaType)
		}
	})

	// Test for method not allowed
	t.Run("Method Not Allowed", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/micropub", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handleMicropub)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
		}
	})
}
