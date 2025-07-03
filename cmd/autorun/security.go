package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Merith-TK/utils/pkg/config"
	"github.com/Merith-TK/utils/pkg/driveutil"
)

// SecurityDecision represents a user's decision about an autorun config
type SecurityDecision int

const (
	SecurityDecisionUnknown SecurityDecision = iota
	SecurityDecisionAllow
	SecurityDecisionAllowOnce
	SecurityDecisionDeny
	SecurityDecisionDenyOnce
)

// String returns a string representation of the security decision
func (s SecurityDecision) String() string {
	switch s {
	case SecurityDecisionAllow:
		return "Allow"
	case SecurityDecisionAllowOnce:
		return "Allow Once"
	case SecurityDecisionDeny:
		return "Deny"
	case SecurityDecisionDenyOnce:
		return "Deny Once"
	default:
		return "Unknown"
	}
}

// ConfigMetadata represents metadata about an autorun config
type ConfigMetadata struct {
	SHA256Hash   string            `json:"sha256_hash"`
	MD5Hash      string            `json:"md5_hash"`
	Decision     SecurityDecision  `json:"decision"`
	LastSeen     time.Time         `json:"last_seen"`
	FirstSeen    time.Time         `json:"first_seen"`
	SeenCount    int               `json:"seen_count"`
	Config       Config            `json:"config"`
	Environment  map[string]string `json:"environment"`
}

// SecurityManager manages security decisions for autorun configs
type SecurityManager struct {
	metadataPath string
	metadata     map[string]*ConfigMetadata
}

// NewSecurityManager creates a new security manager
func NewSecurityManager() *SecurityManager {
	appDataPath := os.Getenv("APPDATA")
	if appDataPath == "" {
		appDataPath = os.Getenv("USERPROFILE")
	}
	
	metadataDir := filepath.Join(appDataPath, "AutorunManager")
	os.MkdirAll(metadataDir, 0755)
	
	sm := &SecurityManager{
		metadataPath: filepath.Join(metadataDir, "security_metadata.json"),
		metadata:     make(map[string]*ConfigMetadata),
	}
	
	sm.loadMetadata()
	return sm
}

// loadMetadata loads existing security metadata from disk
func (sm *SecurityManager) loadMetadata() {
	data, err := os.ReadFile(sm.metadataPath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("[SECURITY] Error reading metadata file: %v\n", err)
		}
		return
	}
	
	if err := json.Unmarshal(data, &sm.metadata); err != nil {
		fmt.Printf("[SECURITY] Error parsing metadata file: %v\n", err)
		return
	}
}

// saveMetadata saves security metadata to disk
func (sm *SecurityManager) saveMetadata() {
	data, err := json.MarshalIndent(sm.metadata, "", "  ")
	if err != nil {
		fmt.Printf("[SECURITY] Error marshaling metadata: %v\n", err)
		return
	}
	
	if err := os.WriteFile(sm.metadataPath, data, 0644); err != nil {
		fmt.Printf("[SECURITY] Error writing metadata file: %v\n", err)
	}
}

// hashConfig creates SHA256 and MD5 hashes for a config
func hashConfig(cfg *Config) (string, string) {
	// Create a deterministic representation of the config
	configData := struct {
		Autorun     string            `json:"autorun"`
		WorkDir     string            `json:"workdir"`
		Isolate     bool              `json:"isolate"`
		Environment map[string]string `json:"environment"`
	}{
		Autorun:     cfg.Autorun,
		WorkDir:     cfg.WorkDir,
		Isolate:     cfg.Isolate,
		Environment: cfg.Environment,
	}
	
	data, _ := json.Marshal(configData)
	
	// Calculate SHA256 hash
	sha256Hash := sha256.Sum256(data)
	sha256Hex := hex.EncodeToString(sha256Hash[:])
	
	// Calculate MD5 hash for display
	md5Hash := md5.Sum(data)
	md5Hex := hex.EncodeToString(md5Hash[:])
	
	return sha256Hex, md5Hex
}

// CheckConfig checks if a config is known and returns the security decision
func (sm *SecurityManager) CheckConfig(drivePath string) (SecurityDecision, *ConfigMetadata, error) {
	configPath := filepath.Join(drivePath, ".autorun.toml")
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return SecurityDecisionAllow, nil, nil // No config file, allow
	}
	
	// Get drive serial number
	driveSerial, err := driveutil.GetVolumeSerialNumber(drivePath)
	if err != nil {
		return SecurityDecisionDeny, nil, fmt.Errorf("failed to get drive serial: %v", err)
	}
	
	// Load the config
	var cfg Config
	if err := config.LoadToml(&cfg, configPath); err != nil {
		return SecurityDecisionDeny, nil, fmt.Errorf("failed to load config: %v", err)
	}
	
	// Calculate hashes
	sha256Hash, md5Hash := hashConfig(&cfg)
	
	// Use drive serial as key
	driveKey := fmt.Sprintf("%08X", driveSerial)
	
	// Check if we have metadata for this drive
	if metadata, exists := sm.metadata[driveKey]; exists {
		// Update last seen and count
		metadata.LastSeen = time.Now()
		metadata.SeenCount++
		sm.saveMetadata()
		
		// Check decision type
		switch metadata.Decision {
		case SecurityDecisionAllow, SecurityDecisionDeny:
			return metadata.Decision, metadata, nil
		case SecurityDecisionAllowOnce, SecurityDecisionDenyOnce:
			// One-time decisions expire after use
			delete(sm.metadata, driveKey)
			sm.saveMetadata()
			return metadata.Decision, metadata, nil
		}
	}
	
	// Create new metadata for unknown config
	metadata := &ConfigMetadata{
		SHA256Hash:  sha256Hash,
		MD5Hash:     md5Hash,
		Decision:    SecurityDecisionUnknown,
		LastSeen:    time.Now(),
		FirstSeen:   time.Now(),
		SeenCount:   1,
		Config:      cfg,
		Environment: cfg.Environment,
	}
	
	return SecurityDecisionUnknown, metadata, nil
}

// SaveDecision saves a security decision for a config
func (sm *SecurityManager) SaveDecision(metadata *ConfigMetadata, decision SecurityDecision, drivePath string) error {
	// Get drive serial number
	driveSerial, err := driveutil.GetVolumeSerialNumber(drivePath)
	if err != nil {
		return fmt.Errorf("failed to get drive serial: %v", err)
	}
	
	// Use drive serial as key
	driveKey := fmt.Sprintf("%08X", driveSerial)
	
	metadata.Decision = decision
	sm.metadata[driveKey] = metadata
	sm.saveMetadata()
	return nil
}

// GetMetadataPath returns the path to the metadata file
func (sm *SecurityManager) GetMetadataPath() string {
	return sm.metadataPath
}

// GetAllMetadata returns all stored metadata
func (sm *SecurityManager) GetAllMetadata() map[string]*ConfigMetadata {
	return sm.metadata
}

// RemoveMetadata removes metadata for a specific drive serial
func (sm *SecurityManager) RemoveMetadata(driveSerial string) {
	delete(sm.metadata, driveSerial)
	sm.saveMetadata()
}

// ClearAllMetadata removes all stored metadata
func (sm *SecurityManager) ClearAllMetadata() {
	sm.metadata = make(map[string]*ConfigMetadata)
	sm.saveMetadata()
}
