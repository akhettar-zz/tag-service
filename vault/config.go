package vault

import (
	"bytes"
	"github.com/BetaProjectWave/kube-vault-plugin"
	"github.com/tag-service/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// Vault config type
var config vault.Config

const (
	// ENVIRONMENT variable set in docker-compose or kubernetes deployment script.
	ENVIRONMENT = "ENVIRONMENT"

	// CONFIG folder path environment variable for this application
	CONFIG = "CONFIG_FOLDER"
)

// LoadConfig loads vault config
func LoadConfig() vault.Config {

	// Build the path to app config
	var pathBuilder bytes.Buffer
	pathBuilder.WriteString(GetEnv(CONFIG, "config"))
	pathBuilder.WriteString("/vault-config-")
	pathBuilder.WriteString(GetEnv(ENVIRONMENT, "default"))
	pathBuilder.WriteString(".yml")

	// Load the config
	logger.Info.Printf("Loading configuration file from: %s", pathBuilder.String())
	source, err := ioutil.ReadFile(pathBuilder.String())
	if err != nil {
		logger.Error.Printf("Failed to load the app config")
		panic(err)
	}
	var config vault.Config
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		logger.Error.Printf("Failed to unmarshal the vault config file")
		panic(err)
	}
	return config
}

// GetEnv env variable or fall back to default
func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
