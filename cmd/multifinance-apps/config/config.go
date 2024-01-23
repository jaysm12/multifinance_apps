package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config struct to hold the configuration data for server
type Config struct {
	Port           string   `yaml:"port"`
	Mysql          Mysql    `yaml:"mysql"`
	Hash           Hash     `yaml:"hash"`
	Token          Token    `yaml:"token"`
	UserHandler    Handler  `yaml:"user_handler"`
	AuthHandler    Handler  `yaml:"auth_handler"`
	PartnerHandler Handler  `yaml:"partner_handler"`
	MaxCounter     int      `yaml:"max_find_counter"`
	RabbitMQ       RabbitMQ `yaml:"rabbit_mq"`
}

// Mysql struct to hold the configuration data for mysql
type Mysql struct {
	Config string `yaml:"mysql_config"`
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

type RabbitMQ struct {
	Config string `yaml:"rabbit_mq_config"`
}

// GetConfig is func to load config and replace it by secret value
func GetConfig() (Config, error) {
	var cfg Config
	// Read the YAML file into a byte slice
	configPath := filepath.Join("config", "config.yaml")
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}

	// Unmarshal the YAML into a Config struct
	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
