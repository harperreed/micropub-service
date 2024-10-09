package main

import (
	"errors"
	"github.com/harperreed/micropub-service/internal/errors"
	"sync"
	"testing"
	"time"
)

// setupTest creates a new MockPocketBase instance for testing
func setupTest() *MockPocketBase {
	return &MockPocketBase{
		metadata: make(map[string]BlogEntryMetadata),
		mu:       sync.RWMutex{},
	}
}

func TestStoreBlogEntryMetadata(t *testing.T) {
	mockPB := setupTest()

	// Test case 1: Successful metadata storage
	t.Run("Successful metadata storage", func(t *testing.T) {
		metadata := BlogEntryMetadata{
			ID:       "1",
			Title:    "Test Post",
			Author:   "Captain Codebeard",
			Date:     "2024-09-05",
			Tags:     []string{"pirate", "coding"},
			Category: "adventures",
		}

		err := StoreBlogEntryMetadata(mockPB, metadata)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	// Test case 2: Failed metadata storage
	t.Run("Failed metadata storage", func(t *testing.T) {
		metadata := BlogEntryMetadata{
			ID:    "2",
			Title: "Invalid Post",
		}

		err := StoreBlogEntryMetadata(mockPB, metadata)
		if err == nil {
			t.Error("Expected an error, got nil")
		}
	})
}

// MockPocketBase is a mock implementation of the PocketBase client
type MockPocketBase struct {
	metadata map[string]BlogEntryMetadata
	mu       sync.RWMutex
}

// BlogEntryMetadata represents the metadata for a blog entry
type BlogEntryMetadata struct {
	ID       string
	Title    string
	Author   string
	Date     string
	Tags     []string
	Category string
}

// StoreBlogEntryMetadata stores blog entry metadata
func StoreBlogEntryMetadata(pb *MockPocketBase, metadata BlogEntryMetadata) error {
	if metadata.ID == "" {
		return errors.NewInvalidMetadataError("ID is required")
	}
	if metadata.Title == "" {
		return errors.NewInvalidMetadataError("Title is required")
	}
	if metadata.Date != "" {
		_, err := time.Parse("2006-01-02", metadata.Date)
		if err != nil {
			return errors.NewInvalidMetadataError("Invalid date format")
		}
	}

	pb.mu.Lock()
	defer pb.mu.Unlock()

	pb.metadata[metadata.ID] = metadata
	return nil
}
