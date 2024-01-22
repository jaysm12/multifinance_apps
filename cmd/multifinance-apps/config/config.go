package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config struct to hold the configuration data for server
type Config struct {
	Port           string  `yaml:"port"`
	Mysql          Mysql   `yaml:"Mysql"`
	Redis          Redis   `yaml:"redis"`
	Hash           Hash    `yaml:"hash"`
	Token          Token   `yaml:"token"`
	UserHandler    Handler `yaml:"user_handler"`
	AuthHandler    Handler `yaml:"auth_handler"`
	PartnerHandler Handler `yaml:"partner_handler"`
	MaxCounter     int     `yaml:"max_find_counter"`
}

// Mysql struct to hold the configuration data for mysql
type Mysql struct {
	Config string `yaml:"postgres_config"`
}

// Redis struct to hold the configuration data for redis
type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
}

// Hash struct to hold the configuration data for Hash Package
type Hash struct {
	Cost int `yaml:"cost"`
}

// Token struct to hold the configuration data for Token Package
type Token struct {
	Secret    string `yaml:"secret"`
	ExpInHour int64  `yaml:"exp_in_hour"`
}

// Handler struct to hold the configuration data for handler
type Handler struct {
	TimeoutInSec int `yaml:"timeout_in_sec"`
}

// GetConfig is func to load config and replace it by secret value
func GetConfig(values map[string]string) (Config, error) {
	var cfg Config
	// Read the YAML file into a byte slice
	configPath := filepath.Join("config", "config.yaml")
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}
	// Replace yaml value with secret
	for k, v := range values {
		configData = []byte(strings.Replace(string(configData), fmt.Sprintf("<%v>", k), v, -1))
	}
	// Unmarshal the YAML into a Config struct
	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
