// config/config.go

package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	configObj *ConfigObject
	cfgPath   string
)

// GetConfig returns the loaded config
func GetConfig() *ConfigObject {
	if configObj == nil {
		var err error
		configObj, err = LoadConfig(cfgPath)
		if err != nil {
			return nil
		}
		return configObj
	}
	return configObj
}

// LoadConfig reads the config file and loads it into ConfigObj
func LoadConfig(configPath string) (*ConfigObject, error) {
	cfgPath = configPath
	file, err := os.Open(configPath)
	if err != nil {
		log.Printf("Failed to open config file: %v", err)
		return nil, err
	}
	defer file.Close()

	var cfg ConfigObject
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		log.Printf("Failed to decode config: %v", err)
		return nil, err
	}
	configObj = &cfg
	return configObj, nil
}
