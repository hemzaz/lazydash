package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	// Create a temporary test file
	testContent := []byte("test file content")
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(testContent); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	// Successful read
	t.Run("Valid file", func(t *testing.T) {
		data, err := LoadFromFile(tmpfile.Name())
		if err != nil {
			t.Fatalf("LoadFromFile(%q) returned error: %v", tmpfile.Name(), err)
		}
		if string(data) != string(testContent) {
			t.Errorf("LoadFromFile(%q) = %q; want %q", tmpfile.Name(), string(data), string(testContent))
		}
	})

	// Non-existent file
	t.Run("Non-existent file", func(t *testing.T) {
		nonExistentPath := filepath.Join(os.TempDir(), "non_existent_file.txt")
		_, err := LoadFromFile(nonExistentPath)
		if err == nil {
			t.Errorf("LoadFromFile with non-existent file should return an error")
		}
	})

	// Directory instead of file
	t.Run("Directory instead of file", func(t *testing.T) {
		_, err := LoadFromFile(os.TempDir())
		if err == nil {
			t.Errorf("LoadFromFile with directory path should return an error")
		}
	})
}

// TestLoadFromStdin is harder to test as it relies on os.Stdin
// In a real-world scenario, you might use dependency injection to make this more testable
// For simplicity here, we'll just make sure the function exists and doesn't panic
func TestLoadFromStdin(t *testing.T) {
	// We can't easily test stdin in a unit test
	// This is just a basic "doesn't crash" test
	t.Run("Function exists", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("LoadFromStdin() panicked: %v", r)
			}
		}()
		
		// Call should return nil, nil in a test environment since stdin is not a pipe
		_, err := LoadFromStdin()
		if err != nil {
			t.Logf("Note: LoadFromStdin() returned error: %v", err)
		}
	})
}