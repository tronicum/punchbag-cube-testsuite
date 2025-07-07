package client

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Try to load a config value from (in order):
// 1. CLI flag (if provided)
// 2. Environment variable
// 3. Config file (YAML or plain text)
func LoadHetznerAPIToken(flagVal string) (string, error) {
	if flagVal != "" {
		return flagVal, nil
	}
	if env := os.Getenv("HETZNER_API_TOKEN"); env != "" {
		return env, nil
	}
	// Look for a config file in $HOME/.hetzner or $PWD/hetzner_token
	home, _ := os.UserHomeDir()
	paths := []string{
		"./hetzner_token",
		filepath.Join(home, ".hetzner_token"),
		filepath.Join(home, ".hetzner", "token"),
	}
	for _, p := range paths {
		if data, err := ioutil.ReadFile(p); err == nil {
			tok := strings.TrimSpace(string(data))
			if tok != "" {
				return tok, nil
			}
		}
	}
	return "", os.ErrNotExist
}

func LoadHetznerS3Credentials() (string, string) {
	accessKey := os.Getenv("HETZNER_S3_ACCESS_KEY")
	secretKey := os.Getenv("HETZNER_S3_SECRET_KEY")
	return accessKey, secretKey
}
