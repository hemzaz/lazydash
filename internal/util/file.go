package util

import (
	"io"
	"os"

	"github.com/rs/zerolog/log"
)

// LoadFromFile loads data from a file
func LoadFromFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// LoadFromStdin loads data from stdin
func LoadFromStdin() ([]byte, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}
	
	// Check if it's a pipe or redirect
	if info.Mode()&os.ModeNamedPipe == 0 && (info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0) {
		return nil, nil
	}

	// Read from stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	
	return data, nil
}