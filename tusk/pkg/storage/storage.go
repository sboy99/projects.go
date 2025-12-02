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
type Storage struct {
	filePath string
	mu       sync.RWMutex
	data     map[string]any
}

// NewStorage creates a new storage instance with the given file path
// The file path can be a filename (stored in current directory) or a full path
func NewStorage(filePath string) *Storage {
	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if dir != "." && dir != "" {
		_ = os.MkdirAll(dir, 0755)
	}

	storage := &Storage{
		filePath: filePath,
		data:     make(map[string]any),
	}

	// Load existing data if file exists (ignore errors, start with empty data)
	_ = storage.load()

	return storage
}

// load reads data from the JSON file
func (s *Storage) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If file doesn't exist, start with empty data
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		s.data = make(map[string]any)
		return nil
	}

	file, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// If file is empty, start with empty data
	if len(file) == 0 {
		s.data = make(map[string]any)
		return nil
	}

	if err := json.Unmarshal(file, &s.data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// save writes data to the JSON file
// Note: This method should only be called while holding a write lock (Lock())
func (s *Storage) save() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Create adds a new record with the given key and value
func (s *Storage) Create(key string, value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; exists {
		return fmt.Errorf("key '%s' already exists", key)
	}

	s.data[key] = value
	return s.save()
}

// Read retrieves a record by key
func (s *Storage) Read(key string) (any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.data[key]
	if !exists {
		return nil, fmt.Errorf("key '%s' not found", key)
	}

	return value, nil
}

// ReadAll retrieves all records
func (s *Storage) ReadAll() (map[string]any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a copy to prevent external modification
	result := make(map[string]any)
	maps.Copy(result, s.data)

	return result, nil
}

// Update updates an existing record
func (s *Storage) Update(key string, value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; !exists {
		return fmt.Errorf("key '%s' not found", key)
	}

	s.data[key] = value
	return s.save()
}

// Upsert creates or updates a record (insert if not exists, update if exists)
func (s *Storage) Upsert(key string, value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	return s.save()
}

// Delete removes a record by key
func (s *Storage) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; !exists {
		return fmt.Errorf("key '%s' not found", key)
	}

	delete(s.data, key)
	return s.save()
}

// Exists checks if a key exists
func (s *Storage) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.data[key]
	return exists
}

// Count returns the number of records
func (s *Storage) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

// Clear removes all records
func (s *Storage) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]any)
	return s.save()
}

// GetFilePath returns the file path used by this storage instance
func (s *Storage) GetFilePath() string {
	return s.filePath
}
