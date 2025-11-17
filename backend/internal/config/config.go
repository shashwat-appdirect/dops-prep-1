package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	FirebaseServiceAccount map[string]interface{}
	SubcollectionID        string
	AdminPassword          string
	Port                   string
	CORSOrigin            string
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Firebase Service Account
	serviceAccount := os.Getenv("FIREBASE_SERVICE_ACCOUNT")
	if serviceAccount == "" {
		return nil, fmt.Errorf("FIREBASE_SERVICE_ACCOUNT environment variable is required")
	}

	var serviceAccountJSON map[string]interface{}

	// Check if it's a base64 encoded string
	if len(serviceAccount) > 7 && serviceAccount[:7] == "base64:" {
		decoded, err := base64.StdEncoding.DecodeString(serviceAccount[7:])
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 service account: %v", err)
		}
		if err := json.Unmarshal(decoded, &serviceAccountJSON); err != nil {
			return nil, fmt.Errorf("failed to parse service account JSON: %v", err)
		}
	} else {
		// Assume it's a file path
		data, err := os.ReadFile(serviceAccount)
		if err != nil {
			return nil, fmt.Errorf("failed to read service account file: %v", err)
		}
		if err := json.Unmarshal(data, &serviceAccountJSON); err != nil {
			return nil, fmt.Errorf("failed to parse service account JSON: %v", err)
		}
	}

	cfg.FirebaseServiceAccount = serviceAccountJSON

	// Subcollection ID
	cfg.SubcollectionID = os.Getenv("SUBSCOLLECTION_ID")
	if cfg.SubcollectionID == "" {
		return nil, fmt.Errorf("SUBSCOLLECTION_ID environment variable is required")
	}

	// Admin Password
	cfg.AdminPassword = os.Getenv("ADMIN_PASSWORD")
	if cfg.AdminPassword == "" {
		return nil, fmt.Errorf("ADMIN_PASSWORD environment variable is required")
	}

	// Port
	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	// CORS Origin
	cfg.CORSOrigin = os.Getenv("CORS_ORIGIN")
	if cfg.CORSOrigin == "" {
		cfg.CORSOrigin = "http://localhost:5173"
	}

	return cfg, nil
}

