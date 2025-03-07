package grafana

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hemzaz/lazydash/internal/config"
	"github.com/hemzaz/lazydash/internal/util"
	"github.com/rs/zerolog/log"
)

// GrafanaFolder represents a folder in Grafana
type GrafanaFolder struct {
	ID        int    `json:"id,omitempty"`
	UID       string `json:"uid,omitempty"`
	Title     string `json:"title"`
	URL       string `json:"url,omitempty"`
	HasACL    bool   `json:"hasAcl,omitempty"`
	CanSave   bool   `json:"canSave,omitempty"`
	CanEdit   bool   `json:"canEdit,omitempty"`
	CanAdmin  bool   `json:"canAdmin,omitempty"`
	CreatedBy string `json:"createdBy,omitempty"`
	Created   string `json:"created,omitempty"`
	UpdatedBy string `json:"updatedBy,omitempty"`
	Updated   string `json:"updated,omitempty"`
	Version   int    `json:"version,omitempty"`
}

// GrafanaFolderError represents an error response from Grafana
type GrafanaFolderError struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// PostDashboard posts a dashboard to Grafana
func PostDashboard(host string, insecureSkipVerify bool, token string, dashboard *Dashboard, cfg *config.Config) {
	if !util.IsURL(host) {
		log.Fatal().Str("host", host).Msg("Host URL is not valid")
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a custom client with optional TLS config
	client := &http.Client{}
	if insecureSkipVerify {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client.Transport = customTransport
		log.Warn().Msg("TLS certificate verification disabled. This is insecure!")
	}
	
	// If folder config is provided, set up dashboard in the specified folder
	if cfg != nil && cfg.FolderConfig != nil && cfg.FolderConfig.Name != "" {
		// Get or create the folder
		folder, err := GetOrCreateFolder(host, token, insecureSkipVerify, cfg.FolderConfig)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to get/create folder. Using default folder")
		} else {
			// Set the folder ID in the dashboard
			dashboard.FolderID = folder.ID
			log.Info().Str("folder", folder.Title).Int("id", folder.ID).Msg("Using folder")
		}
	}

	// Prepare dashboard submission
	dashboardSubmission := struct {
		Dashboard *Dashboard `json:"dashboard"`
		FolderID  int        `json:"folderId,omitempty"`
		Overwrite bool       `json:"overwrite"`
	}{
		Dashboard: dashboard,
		FolderID:  dashboard.FolderID,
		Overwrite: true,
	}

	// Marshal to JSON
	body, err := json.Marshal(dashboardSubmission)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create JSON body")
	}
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, host+"/api/dashboards/db", bytes.NewReader(body))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create request")
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to post dashboard")
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		log.Fatal().Int("status", resp.StatusCode).Str("body", string(responseBody)).Msg("HTTP request failed")
	}
	
	// Read response
	var dashboardResp struct {
		ID      int    `json:"id"`
		UID     string `json:"uid"`
		URL     string `json:"url"`
		Status  string `json:"status"`
		Version int    `json:"version"`
	}
	
	responseBody, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(responseBody, &dashboardResp); err == nil && dashboardResp.URL != "" {
		log.Info().Str("dashboard", dashboard.Title).Str("url", host+dashboardResp.URL).Msg("Dashboard created")
	} else {
		log.Info().Str("dashboard", dashboard.Title).Msg("Dashboard created successfully")
	}
}

// GetOrCreateFolder gets a folder by title, or creates it if it doesn't exist
func GetOrCreateFolder(host, token string, insecureSkipVerify bool, config *config.FolderConfig) (*GrafanaFolder, error) {
	// First try to find existing folder
	folders, err := GetFoldersByTitle(host, token, insecureSkipVerify, config.Name)
	if err != nil {
		return nil, fmt.Errorf("error getting folders: %w", err)
	}
	
	// Return existing folder if found
	if len(folders) > 0 {
		return &folders[0], nil
	}
	
	// Create new folder if allowed
	if config.Create {
		folder := &GrafanaFolder{
			Title: config.Name,
		}
		return CreateFolder(host, token, insecureSkipVerify, folder)
	}
	
	return nil, fmt.Errorf("folder '%s' not found and create=false", config.Name)
}

// GetFoldersByTitle gets folders from Grafana API filtered by title
func GetFoldersByTitle(host, token string, insecureSkipVerify bool, title string) ([]GrafanaFolder, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a custom client with optional TLS config
	client := &http.Client{}
	if insecureSkipVerify {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client.Transport = customTransport
	}
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, host+"/api/folders", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	
	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting folders: %w", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}
	
	// Read response
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	
	// Parse response
	var folders []GrafanaFolder
	if err := json.Unmarshal(responseBody, &folders); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}
	
	// Filter by title if provided
	if title != "" {
		var filtered []GrafanaFolder
		for _, folder := range folders {
			if folder.Title == title {
				filtered = append(filtered, folder)
			}
		}
		return filtered, nil
	}
	
	return folders, nil
}

// CreateFolder creates a new folder in Grafana
func CreateFolder(host, token string, insecureSkipVerify bool, folder *GrafanaFolder) (*GrafanaFolder, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a custom client with optional TLS config
	client := &http.Client{}
	if insecureSkipVerify {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client.Transport = customTransport
	}
	
	// Marshal folder to JSON
	body, err := json.Marshal(folder)
	if err != nil {
		return nil, fmt.Errorf("error marshaling folder: %w", err)
	}
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, host+"/api/folders", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	
	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error creating folder: %w", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}
	
	// Read response
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	
	// Parse response
	var createdFolder GrafanaFolder
	if err := json.Unmarshal(responseBody, &createdFolder); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}
	
	return &createdFolder, nil
}