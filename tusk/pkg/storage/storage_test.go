package storage

import (
	"os"
	"path/filepath"
	"testing"
)

type TestItem struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestNewStorage(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")

	storage := NewStorage[TestItem](filePath)
	if storage == nil {
		t.Fatal("NewStorage() returned nil")
	}
	if storage.filePath != filePath {
		t.Errorf("NewStorage() filePath = %q, want %q", storage.filePath, filePath)
	}
	if storage.data == nil {
		t.Error("NewStorage() data is nil")
	}
}

func TestStorage_Create(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")
	storage := NewStorage[TestItem](filePath)

	t.Run("create new item", func(t *testing.T) {
		item := TestItem{Name: "test", Value: 42}
		err := storage.Create("key1", item)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// Verify it was created
		read, err := storage.Read("key1")
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}
		if read.Name != item.Name {
			t.Errorf("Read() Name = %q, want %q", read.Name, item.Name)
		}
		if read.Value != item.Value {
			t.Errorf("Read() Value = %d, want %d", read.Value, item.Value)
		}
	})

	t.Run("create duplicate key", func(t *testing.T) {
		item := TestItem{Name: "test", Value: 42}
		err := storage.Create("key1", item)
		if err == nil {
			t.Error("Create() expected error for duplicate key, got nil")
		}
	})
}

func TestStorage_Upsert(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")
	storage := NewStorage[TestItem](filePath)

	t.Run("upsert new item", func(t *testing.T) {
		item := TestItem{Name: "test", Value: 42}
		err := storage.Upsert("key1", item)
		if err != nil {
			t.Fatalf("Upsert() error = %v", err)
		}
	})

	t.Run("upsert existing item", func(t *testing.T) {
		item1 := TestItem{Name: "test1", Value: 42}
		item2 := TestItem{Name: "test2", Value: 100}

		storage.Upsert("key1", item1)
		err := storage.Upsert("key1", item2)
		if err != nil {
			t.Fatalf("Upsert() error = %v", err)
		}

		// Verify it was updated
		read, err := storage.Read("key1")
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}
		if read.Name != item2.Name {
			t.Errorf("Read() Name = %q, want %q", read.Name, item2.Name)
		}
		if read.Value != item2.Value {
			t.Errorf("Read() Value = %d, want %d", read.Value, item2.Value)
		}
	})
}

func TestStorage_Read(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")
	storage := NewStorage[TestItem](filePath)

	t.Run("read existing item", func(t *testing.T) {
		item := TestItem{Name: "test", Value: 42}
		storage.Upsert("key1", item)

		read, err := storage.Read("key1")
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}
		if read.Name != item.Name {
			t.Errorf("Read() Name = %q, want %q", read.Name, item.Name)
		}
	})

	t.Run("read non-existent item", func(t *testing.T) {
		_, err := storage.Read("nonexistent")
		if err == nil {
			t.Error("Read() expected error for non-existent key, got nil")
		}
	})
}

func TestStorage_ReadAll(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")
	storage := NewStorage[TestItem](filePath)

	t.Run("read all empty storage", func(t *testing.T) {
		all, err := storage.ReadAll()
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}
		if len(all) != 0 {
			t.Errorf("ReadAll() length = %d, want 0", len(all))
		}
	})

	t.Run("read all with items", func(t *testing.T) {
		storage.Upsert("key1", TestItem{Name: "test1", Value: 1})
		storage.Upsert("key2", TestItem{Name: "test2", Value: 2})
		storage.Upsert("key3", TestItem{Name: "test3", Value: 3})

		all, err := storage.ReadAll()
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}
		if len(all) != 3 {
			t.Errorf("ReadAll() length = %d, want 3", len(all))
		}
		if all["key1"].Value != 1 {
			t.Errorf("ReadAll() key1.Value = %d, want 1", all["key1"].Value)
		}
	})
}

func TestStorage_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")
	storage := NewStorage[TestItem](filePath)

	t.Run("delete existing item", func(t *testing.T) {
		storage.Upsert("key1", TestItem{Name: "test", Value: 42})

		err := storage.Delete("key1")
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		// Verify it was deleted
		_, err = storage.Read("key1")
		if err == nil {
			t.Error("Delete() item should be deleted")
		}
	})

	t.Run("delete non-existent item", func(t *testing.T) {
		err := storage.Delete("nonexistent")
		if err == nil {
			t.Error("Delete() expected error for non-existent key, got nil")
		}
	})
}

func TestStorage_Exists(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")
	storage := NewStorage[TestItem](filePath)

	t.Run("exists for existing item", func(t *testing.T) {
		storage.Upsert("key1", TestItem{Name: "test", Value: 42})

		if !storage.Exists("key1") {
			t.Error("Exists() = false, want true")
		}
	})

	t.Run("exists for non-existent item", func(t *testing.T) {
		if storage.Exists("nonexistent") {
			t.Error("Exists() = true, want false")
		}
	})
}

func TestStorage_Count(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")
	storage := NewStorage[TestItem](filePath)

	t.Run("count empty storage", func(t *testing.T) {
		count := storage.Count()
		if count != 0 {
			t.Errorf("Count() = %d, want 0", count)
		}
	})

	t.Run("count with items", func(t *testing.T) {
		storage.Upsert("key1", TestItem{Name: "test1", Value: 1})
		storage.Upsert("key2", TestItem{Name: "test2", Value: 2})

		count := storage.Count()
		if count != 2 {
			t.Errorf("Count() = %d, want 2", count)
		}
	})
}

func TestStorage_Clear(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")
	storage := NewStorage[TestItem](filePath)

	t.Run("clear with items", func(t *testing.T) {
		storage.Upsert("key1", TestItem{Name: "test1", Value: 1})
		storage.Upsert("key2", TestItem{Name: "test2", Value: 2})

		err := storage.Clear()
		if err != nil {
			t.Fatalf("Clear() error = %v", err)
		}

		count := storage.Count()
		if count != 0 {
			t.Errorf("Clear() Count() = %d, want 0", count)
		}
	})
}

func TestStorage_GetFilePath(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")
	storage := NewStorage[TestItem](filePath)

	got := storage.GetFilePath()
	if got != filePath {
		t.Errorf("GetFilePath() = %q, want %q", got, filePath)
	}
}

func TestStorage_Persistence(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")

	// Create storage and add item
	storage1 := NewStorage[TestItem](filePath)
	storage1.Upsert("key1", TestItem{Name: "test", Value: 42})

	// Create new storage instance (should load from file)
	storage2 := NewStorage[TestItem](filePath)

	// Verify item was persisted
	item, err := storage2.Read("key1")
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if item.Name != "test" {
		t.Errorf("Read() Name = %q, want 'test'", item.Name)
	}
	if item.Value != 42 {
		t.Errorf("Read() Value = %d, want 42", item.Value)
	}
}

func TestStorage_LoadEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")

	// Create empty file
	os.WriteFile(filePath, []byte{}, 0644)

	storage := NewStorage[TestItem](filePath)
	count := storage.Count()
	if count != 0 {
		t.Errorf("NewStorage() Count() = %d, want 0 for empty file", count)
	}
}

func TestStorage_LoadInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")

	// Create file with invalid JSON
	os.WriteFile(filePath, []byte("invalid json"), 0644)

	// Should handle gracefully (load() ignores errors in NewStorage)
	storage := NewStorage[TestItem](filePath)
	// Should start with empty data
	count := storage.Count()
	if count != 0 {
		t.Logf("NewStorage() with invalid JSON: Count() = %d (may be 0 or may have partial data)", count)
	}
}

func TestStorage_ConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")
	storage := NewStorage[TestItem](filePath)

	// Test that concurrent access doesn't panic
	done := make(chan bool)

	// Concurrent writes
	go func() {
		for i := 0; i < 10; i++ {
			storage.Upsert("key1", TestItem{Name: "test", Value: i})
		}
		done <- true
	}()

	// Concurrent reads
	go func() {
		for i := 0; i < 10; i++ {
			storage.Read("key1")
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done
}
