package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const (
	maxUploadSize = 10 << 20 // 10 MB
	uploadsDir    = "uploads"
)

var allowedFileTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

func handleMediaUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit the request body size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check file type
	if !isAllowedFileType(file) {
		http.Error(w, "File type not allowed", http.StatusBadRequest)
		return
	}

	// Create the uploads directory if it doesn't exist
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create uploads directory", http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	dst, err := os.Create(filepath.Join(uploadsDir, handler.Filename))
	if err != nil {
		http.Error(w, "Failed to create the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save the file", http.StatusInternalServerError)
		return
	}

	// Return the URL of the uploaded file
	fileURL := fmt.Sprintf("/%s/%s", uploadsDir, handler.Filename)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"url": "%s"}`, fileURL)
}

func isAllowedFileType(file multipart.File) bool {
	// Read the first 512 bytes to determine the content type
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return false
	}

	// Seek back to the start of the file
	file.Seek(0, 0)

	// Get the content type and check if it's allowed
	contentType := http.DetectContentType(buffer)
	return allowedFileTypes[contentType]
}

func main() {
	http.HandleFunc("/media", handleMediaUpload)
	fmt.Println("Media endpoint server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
