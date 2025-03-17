package util

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid URL with HTTPS",
			input:    "https://example.com",
			expected: true,
		},
		{
			name:     "Valid URL with HTTP",
			input:    "http://example.com",
			expected: true,
		},
		{
			name:     "Valid URL with path",
			input:    "https://example.com/path",
			expected: true,
		},
		{
			name:     "Valid URL with query params",
			input:    "https://example.com?param=value",
			expected: true,
		},
		{
			name:     "Invalid URL - missing scheme",
			input:    "example.com",
			expected: false,
		},
		{
			name:     "Invalid URL - empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "Invalid URL - malformed",
			input:    "https://",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsURL(tt.input)
			if result != tt.expected {
				t.Errorf("IsURL(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFetchURL(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test data"))
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Success case
	t.Run("Success", func(t *testing.T) {
		url := server.URL + "/success"
		data, err := FetchURL(url, false)
		if err != nil {
			t.Fatalf("FetchURL(%q) returned error: %v", url, err)
		}
		if string(data) != "test data" {
			t.Errorf("FetchURL(%q) = %q; want %q", url, string(data), "test data")
		}
	})

	// Invalid URL
	t.Run("Invalid URL", func(t *testing.T) {
		_, err := FetchURL("invalid-url", false)
		if err == nil {
			t.Errorf("FetchURL with invalid URL should return an error")
		}
	})

	// HTTP error status code
	t.Run("HTTP Error", func(t *testing.T) {
		url := server.URL + "/error"
		_, err := FetchURL(url, false)
		if err == nil {
			t.Errorf("FetchURL(%q) should return an error for non-2xx status code", url)
		}
	})

	// Non-existent URL
	t.Run("Non-existent URL", func(t *testing.T) {
		_, err := FetchURL("https://non-existent-url-123456789.example", false)
		if err == nil {
			t.Errorf("FetchURL with non-existent URL should return an error")
		}
	})
}