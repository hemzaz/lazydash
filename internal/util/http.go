// Package util provides common utility functions for lazydash
package util

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
)

// IsURL checks if a string is a valid URL
func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// FetchURL fetches a HTTP URL with a timeout and returns the body as bytes
func FetchURL(urlStr string, insecureSkipVerify bool) []byte {
	if !IsURL(urlStr) {
		log.Fatal().Str("url", urlStr).Msg("URL is not valid")
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Create a custom client if needed
	client := &http.Client{}
	if insecureSkipVerify {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client.Transport = customTransport
		log.Warn().Msg("TLS certificate verification disabled. This is insecure!")
	}
	
	// Create request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		log.Fatal().Err(err).Str("url", urlStr).Msg("Failed to create request")
	}
	
	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal().Err(err).Str("url", urlStr).Msg("Failed to connect to url")
	}
	
	// Always close body when done
	defer resp.Body.Close()
	
	// Check HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Fatal().Int("status", resp.StatusCode).Str("url", urlStr).Msg("HTTP request failed")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Err(err).Str("url", urlStr).Msg("Failed to read response body")
	}

	return body
}