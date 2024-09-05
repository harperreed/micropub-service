package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePost(t *testing.T) {
	// Test for x-www-form-urlencoded request
	t.Run("x-www-form-urlencoded", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/micropub", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handleMicropub)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotImplemented {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotImplemented)
		}

		expected := "Micropub endpoint not yet implemented"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	// Test for JSON request
	t.Run("JSON", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/micropub", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handleMicropub)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotImplemented {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotImplemented)
		}

		expected := "Micropub endpoint not yet implemented"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
}
