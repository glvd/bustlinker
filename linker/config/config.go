package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type CacheConfig struct {
	BackupSeconds int
}

type Pinning struct {
	PerSeconds int
}

type Config struct {
	MaxAttempts int64
	Pinning     Pinning
	Hash        CacheConfig
	Address     CacheConfig
}

//var DefaultBootstrapAddresses = []string{}
var DefaultPinningSeconds = 30
var DefaultConfigName = "linker"

// Clone copies the config. Use when updating.
func (c *Config) Clone() (*Config, error) {
	var newConfig Config
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(c); err != nil {
		return nil, fmt.Errorf("failure to encode config: %s", err)
	}

	if err := json.NewDecoder(&buf).Decode(&newConfig); err != nil {
		return nil, fmt.Errorf("failure to decode config: %s", err)
	}

	return &newConfig, nil
}

func FromMap(v map[string]interface{}) (*Config, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}
	var conf Config
	if err := json.NewDecoder(buf).Decode(&conf); err != nil {
		return nil, fmt.Errorf("failure to decode config: %s", err)
	}
	return &conf, nil
}

func ToMap(cfg *Config) (map[string]interface{}, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(cfg); err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.NewDecoder(buf).Decode(&m); err != nil {
		return nil, fmt.Errorf("failure to decode config: %s", err)
	}
	return m, nil
}

func StoreConfig(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(path, DefaultConfigName), data, 0755)
}

func InitConfig(path string) *Config {
	open, err := os.Open(filepath.Join(path, DefaultConfigName))
	if err != nil {
		return defaultConfig()
	}
	var cfg Config
	cfgData, err := ioutil.ReadAll(open)
	if err != nil {
		return defaultConfig()
	}
	err = json.Unmarshal(cfgData, &cfg)
	if err != nil {
		return defaultConfig()
	}

	return &cfg
}

func defaultConfig() *Config {
	cfg := Config{
		MaxAttempts: 3,
		Pinning: Pinning{
			PerSeconds: DefaultPinningSeconds,
		},
		Hash: CacheConfig{
			BackupSeconds: 30,
		},
		Address: CacheConfig{
			BackupSeconds: 30,
		},
	}
	return &cfg
}
