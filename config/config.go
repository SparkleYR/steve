package config

import (
	"encoding/json"
	"os"
	"strings"
)

type Config struct {
	Sudo []string `json:"sudo"`
}

var c Config

func LoadConfig() error {
	f, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return err
	}
	return nil
}

func GetConfig() *Config {
	return &c
}

func IsSudo(user string) bool {
	for _, sudo := range c.Sudo {
		// Check exact match first (for backward compatibility)
		if sudo == user {
			return true
		}
		// If sudo contains @, extract the User part and compare
		// This handles both old format (phone@s.whatsapp.net) and new format (lid@lid)
		if len(sudo) > 0 {
			// Try to extract just the numeric part before @
			parts := strings.Split(sudo, "@")
			if len(parts) > 0 && parts[0] == user {
				return true
			}
		}
	}
	return false
}
