package storage

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"sync"
)

// Storage provides a JSON-based storage layer with CRUD operations
// It stores typed data as a map[string]T
type Storage[T any] struct {
	filePath string
	mu       sync.RWMutex
	data     map[string]T
}

// NewStorage creates a new storage instance with the given file path
// The file path can be a filename (stored in current directory) or a full path
func NewStorage[T any](filePath string) *Storage[T] {
	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if dir != "." && dir != "" {
		_ = os.MkdirAll(dir, 0755)
	}

	storage := &Storage[T]{
		filePath: filePath,
		data:     make(map[string]T),
	}

	// Load existing data if file exists (ignore errors, start with empty data)
	_ = storage.load()

	return storage
}

// load reads data from the JSON file
func (s *Storage[T]) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If file doesn't exist, start with empty data
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		s.data = make(map[string]T)
		return nil
	}

	file, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// If file is empty, start with empty data
	if len(file) == 0 {
		s.data = make(map[string]T)
		return nil
	}

	if err := json.Unmarshal(file, &s.data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// save writes data to the JSON file
// Note: This method should only be called while holding a write lock (Lock())
func (s *Storage[T]) save() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Create adds a new item with the given key
func (s *Storage[T]) Create(key string, item T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; exists {
		return fmt.Errorf("key '%s' already exists", key)
	}

	s.data[key] = item
	return s.save()
}

// Upsert creates or updates an item with the given key
func (s *Storage[T]) Upsert(key string, item T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = item
	return s.save()
}

// Read retrieves an item by key
func (s *Storage[T]) Read(key string) (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var zero T
	value, exists := s.data[key]
	if !exists {
		return zero, fmt.Errorf("key '%s' not found", key)
	}

	return value, nil
}

// ReadAll retrieves all items from the map
func (s *Storage[T]) ReadAll() (map[string]T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a copy to prevent external modification
	result := make(map[string]T)
	maps.Copy(result, s.data)

	return result, nil
}

// GetFilePath returns the file path used by this storage instance
func (s *Storage[T]) GetFilePath() string {
	return s.filePath
}

// Count returns the number of items in the map
func (s *Storage[T]) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

// Clear removes all items from the map
func (s *Storage[T]) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]T)
	return s.save()
}

// Delete removes an item by key
func (s *Storage[T]) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; !exists {
		return fmt.Errorf("key '%s' not found", key)
	}

	delete(s.data, key)
	return s.save()
}

// Exists checks if a key exists
func (s *Storage[T]) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.data[key]
	return exists
}
