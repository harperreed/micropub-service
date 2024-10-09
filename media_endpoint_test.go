package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestMediaEndpoint(t *testing.T) {
	// Setup: Create a temporary directory for uploads
	tempDir, err := os.MkdirTemp("", "test-uploads")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Replace the global uploadsDir with our temp directory
	originalUploadsDir := uploadsDir
	uploadsDir = tempDir
	defer func() { uploadsDir = originalUploadsDir }()

	// Test case 1: Successful file upload
	t.Run("SuccessfulUpload", func(t *testing.T) {
		body, contentType := createMultipartFormData("file", "test.jpg", []byte("fake image content"))
		req, _ := http.NewRequest("POST", "/media", body)
		req.Header.Set("Content-Type", contentType)
		rr := httptest.NewRecorder()

		handleMediaUpload(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var response map[string]string
		json.Unmarshal(rr.Body.Bytes(), &response)
		if _, exists := response["url"]; !exists {
			t.Errorf("response does not contain 'url' field")
		}
	})

	// Test case 2: Missing file in request
	t.Run("MissingFile", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/media", nil)
		rr := httptest.NewRecorder()

		handleMediaUpload(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	// Test case 3: File size exceeding limit
	t.Run("FileSizeExceedingLimit", func(t *testing.T) {
		largeContent := make([]byte, maxUploadSize+1)
		body, contentType := createMultipartFormData("file", "large.jpg", largeContent)
		req, _ := http.NewRequest("POST", "/media", body)
		req.Header.Set("Content-Type", contentType)
		rr := httptest.NewRecorder()

		handleMediaUpload(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	// Test case 4: Unsupported file type
	t.Run("UnsupportedFileType", func(t *testing.T) {
		body, contentType := createMultipartFormData("file", "test.txt", []byte("text content"))
		req, _ := http.NewRequest("POST", "/media", body)
		req.Header.Set("Content-Type", contentType)
		rr := httptest.NewRecorder()

		handleMediaUpload(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	// Test case 5: Concurrent file uploads
	t.Run("ConcurrentUploads", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				body, contentType := createMultipartFormData("file", fmt.Sprintf("test%d.jpg", i), []byte("fake image content"))
				req, _ := http.NewRequest("POST", "/media", body)
				req.Header.Set("Content-Type", contentType)
				rr := httptest.NewRecorder()

				handleMediaUpload(rr, req)

				if status := rr.Code; status != http.StatusOK {
					t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
				}
			}(i)
		}
		wg.Wait()
	})

	// Test case 6: Error in file storage operations
	t.Run("FileStorageError", func(t *testing.T) {
		// Make the uploads directory read-only to simulate a storage error
		os.Chmod(tempDir, 0555)
		defer os.Chmod(tempDir, 0755)

		body, contentType := createMultipartFormData("file", "test.jpg", []byte("fake image content"))
		req, _ := http.NewRequest("POST", "/media", body)
		req.Header.Set("Content-Type", contentType)
		rr := httptest.NewRecorder()

		handleMediaUpload(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}

func createMultipartFormData(fieldName, fileName string, fileContent []byte) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(fieldName, fileName)
	part.Write(fileContent)
	writer.Close()
	return body, writer.FormDataContentType()
}
