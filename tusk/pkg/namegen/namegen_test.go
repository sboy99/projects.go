package namegen

import (
	"strings"
	"testing"
)

func TestNewNameGenerator(t *testing.T) {
	ng := NewNameGenerator()
	if ng == nil {
		t.Fatal("NewNameGenerator() returned nil")
	}
	if ng.usedNames == nil {
		t.Error("NewNameGenerator() usedNames is nil")
	}
	if ng.rng == nil {
		t.Error("NewNameGenerator() rng is nil")
	}
}

func TestNameGenerator_Generate(t *testing.T) {
	ng := NewNameGenerator()

	t.Run("generate unique names", func(t *testing.T) {
		names := make(map[string]bool)
		for i := 0; i < 10; i++ {
			name := ng.Generate()
			if name == "" {
				t.Error("Generate() returned empty string")
			}
			if names[name] {
				t.Errorf("Generate() returned duplicate name: %q", name)
			}
			names[name] = true

			// Check format: adjective-noun-suffix
			parts := strings.Split(name, "-")
			if len(parts) != 3 {
				t.Errorf("Generate() name format incorrect: %q (expected 3 parts)", name)
			}
		}
	})

	t.Run("generate has correct format", func(t *testing.T) {
		name := ng.Generate()
		parts := strings.Split(name, "-")
		if len(parts) != 3 {
			t.Errorf("Generate() name = %q, want format 'adjective-noun-suffix'", name)
		}
	})
}

func TestGenerateOne(t *testing.T) {
	t.Run("generate one name", func(t *testing.T) {
		name := GenerateOne()
		if name == "" {
			t.Error("GenerateOne() returned empty string")
		}

		// Check format
		parts := strings.Split(name, "-")
		if len(parts) != 3 {
			t.Errorf("GenerateOne() name = %q, want format 'adjective-noun-suffix'", name)
		}
	})

	t.Run("generate multiple unique names", func(t *testing.T) {
		names := make(map[string]bool)
		for i := 0; i < 100; i++ {
			name := GenerateOne()
			if names[name] {
				// With 100 attempts, some duplicates are possible but unlikely
				// If we get many duplicates, there might be an issue
				t.Logf("GenerateOne() returned duplicate name: %q (attempt %d)", name, i+1)
			}
			names[name] = true
		}
	})
}

func TestNameGenerator_generateSuffix(t *testing.T) {
	ng := NewNameGenerator()

	t.Run("generate suffix has correct length", func(t *testing.T) {
		suffix := ng.generateSuffix()
		if len(suffix) < 4 || len(suffix) > 5 {
			t.Errorf("generateSuffix() length = %d, want 4 or 5", len(suffix))
		}
	})

	t.Run("generate multiple suffixes", func(t *testing.T) {
		suffixes := make(map[string]bool)
		for i := 0; i < 20; i++ {
			suffix := ng.generateSuffix()
			if len(suffix) < 4 || len(suffix) > 5 {
				t.Errorf("generateSuffix() length = %d, want 4 or 5", len(suffix))
			}
			suffixes[suffix] = true
		}
		// Most should be unique (though some duplicates are possible)
		if len(suffixes) < 15 {
			t.Logf("generateSuffix() generated only %d unique suffixes out of 20", len(suffixes))
		}
	})
}

func TestAddSuffix(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"add suffix to name", "test-name"},
		{"add suffix to simple name", "test"},
		{"add suffix to empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddSuffix(tt.input)
			if result == "" {
				t.Error("AddSuffix() returned empty string")
			}
			if !strings.HasPrefix(result, tt.input) {
				t.Errorf("AddSuffix(%q) = %q, want to start with %q", tt.input, result, tt.input)
			}
			// Should have the original name plus a dash and suffix
			if tt.input != "" && !strings.Contains(result, tt.input+"-") {
				t.Errorf("AddSuffix(%q) = %q, want to contain %q-", tt.input, result, tt.input)
			}
		})
	}
}

func TestNameGenerator_Uniqueness(t *testing.T) {
	ng := NewNameGenerator()

	// Generate many names and check for reasonable uniqueness
	names := make(map[string]bool)
	duplicates := 0

	for i := 0; i < 1000; i++ {
		name := ng.Generate()
		if names[name] {
			duplicates++
		}
		names[name] = true
	}

	// With 1000 names, we expect very few duplicates (close to 0)
	// But due to random generation, a few duplicates are acceptable
	if duplicates > 10 {
		t.Errorf("Generate() produced %d duplicates out of 1000 names (expected < 10)", duplicates)
	}
}

func TestAdjectivesAndNouns(t *testing.T) {
	// Verify that adjectives and nouns arrays are not empty
	if len(adjectives) == 0 {
		t.Error("adjectives array is empty")
	}
	if len(nouns) == 0 {
		t.Error("nouns array is empty")
	}

	// Verify they contain expected values
	foundAdjective := false
	foundNoun := false

	for _, adj := range adjectives {
		if adj != "" {
			foundAdjective = true
			break
		}
	}

	for _, noun := range nouns {
		if noun != "" {
			foundNoun = true
			break
		}
	}

	if !foundAdjective {
		t.Error("adjectives array contains no non-empty values")
	}
	if !foundNoun {
		t.Error("nouns array contains no non-empty values")
	}
}
