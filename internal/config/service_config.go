package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type ServiceConfig struct {
	Service             ServiceDetails  `mapstructure:"service"`
	Shared              SharedConfig    `mapstructure:"shared"`
	Database            DatabaseConfig  `mapstructure:"database"`
	PlayerDatabase      DatabaseConfig  `mapstructure:"player_database"`
	ClawmachineDatabase DatabaseConfig  `mapstructure:"clawmachine_database"`
	Redis               RedisConfig     `mapstructure:"redis"`
	GRPC                GRPCConfig      `mapstructure:"grpc"`
	WebSocket           WebSocketConfig `mapstructure:"websocket"`
	TCP                 TCPConfig       `mapstructure:"tcp"`
	CORS                CORSConfig      `mapstructure:"cors"`
	Logging             LoggingConfig   `mapstructure:"logging"`
	JWT                 JWTConfig       `mapstructure:"jwt"`
	Tracing             TracingConfig   `mapstructure:"tracing"`
	Discovery           DiscoveryConfig `mapstructure:"discovery"`
}

type ServiceDetails struct {
	Name         string `mapstructure:"name"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	Path         string `mapstructure:"path"`
}

type SharedConfig struct {
	Database            string `mapstructure:"database"`
	PlayerDatabase      string `mapstructure:"player_database"`
	ClawmachineDatabase string `mapstructure:"clawmachine_database"`
	Redis               string `mapstructure:"redis"`
	Logging             string `mapstructure:"logging"`
	JWT                 string `mapstructure:"jwt"`
	Tracing             string `mapstructure:"tracing"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

// LoadServiceConfig loads config for a specific service
func LoadServiceConfig(serviceName string) (*ServiceConfig, error) {
	configFile := fmt.Sprintf("config/%s.yaml", serviceName)
	return LoadServiceConfigFromPath(configFile)
}

// LoadServiceConfigFromPath loads config from a specific file path
func LoadServiceConfigFromPath(configFile string) (*ServiceConfig, error) {
	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configFile)
	}

	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")

	// Set environment variable prefix based on config file name
	fileName := filepath.Base(configFile)
	serviceName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	envPrefix := fmt.Sprintf("GAME_%s", strings.ToUpper(strings.ReplaceAll(serviceName, "-", "_")))
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file %s: %w", configFile, err)
	}

	var config ServiceConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Load shared configs if specified
	if config.Shared.Database != "" {
		if err := loadSharedConfig(v, config.Shared.Database, "database"); err != nil {
			return nil, err
		}
	}
	if config.Shared.PlayerDatabase != "" {
		if err := loadSharedConfig(v, config.Shared.PlayerDatabase, "player_database"); err != nil {
			return nil, err
		}
	}
	if config.Shared.ClawmachineDatabase != "" {
		if err := loadSharedConfig(v, config.Shared.ClawmachineDatabase, "clawmachine_database"); err != nil {
			return nil, err
		}
	}
	if config.Shared.Redis != "" {
		if err := loadSharedConfig(v, config.Shared.Redis, "redis"); err != nil {
			return nil, err
		}
	}
	if config.Shared.Logging != "" {
		if err := loadSharedConfig(v, config.Shared.Logging, "logging"); err != nil {
			return nil, err
		}
	}
	if config.Shared.JWT != "" {
		if err := loadSharedConfig(v, config.Shared.JWT, "jwt"); err != nil {
			return nil, err
		}
	}
	if config.Shared.Tracing != "" {
		if err := loadSharedConfig(v, config.Shared.Tracing, "tracing"); err != nil {
			return nil, err
		}
	}

	// Re-unmarshal after loading shared configs
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling final config: %w", err)
	}

	return &config, nil
}

func loadSharedConfig(v *viper.Viper, configFile, key string) error {
	// Resolve relative path from the main config file's directory
	mainConfigDir := filepath.Dir(v.ConfigFileUsed())
	sharedConfigPath := filepath.Join(mainConfigDir, configFile)

	sharedV := viper.New()
	sharedV.SetConfigFile(sharedConfigPath)
	sharedV.SetConfigType("yaml")

	if err := sharedV.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading shared config %s: %w", sharedConfigPath, err)
	}

	// Merge shared config into main config
	sharedSettings := sharedV.AllSettings()
	for k, val := range sharedSettings {
		if k == key {
			// If key matches section name, merge inner map directly
			if innerMap, ok := val.(map[string]interface{}); ok {
				for innerK, innerV := range innerMap {
					v.Set(fmt.Sprintf("%s.%s", key, innerK), innerV)
				}
			}
		} else {
			// Otherwise, set it as is
			v.Set(fmt.Sprintf("%s.%s", key, k), val)
		}
	}

	return nil
}

// Helper methods for getting service addresses
func (c *ServiceConfig) GetServiceAddr() string {
	return fmt.Sprintf("%s:%d", c.Service.Host, c.Service.Port)
}

func (c *ServiceConfig) GetGRPCAddr() string {
	if c.GRPC.Port != 0 {
		return fmt.Sprintf("%s:%d", c.GRPC.Host, c.GRPC.Port)
	}
	return c.GetServiceAddr()
}

func (c *ServiceConfig) GetWebSocketAddr() string {
	if c.WebSocket.Port != 0 {
		return fmt.Sprintf("%s:%d", c.WebSocket.Host, c.WebSocket.Port)
	}
	return c.GetServiceAddr()
}

func (c *ServiceConfig) GetTCPAddr() string {
	if c.TCP.Port != 0 {
		return fmt.Sprintf("%s:%d", c.TCP.Host, c.TCP.Port)
	}
	return c.GetServiceAddr()
}

// GetDSN returns database connection string
func (c *ServiceConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)
}

// GetPlayerDSN returns player database connection string
func (c *ServiceConfig) GetPlayerDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		c.PlayerDatabase.User,
		c.PlayerDatabase.Password,
		c.PlayerDatabase.Host,
		c.PlayerDatabase.Port,
		c.PlayerDatabase.Name,
	)
}

// GetClawmachineDSN returns clawmachine database connection string
func (c *ServiceConfig) GetClawmachineDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		c.ClawmachineDatabase.User,
		c.ClawmachineDatabase.Password,
		c.ClawmachineDatabase.Host,
		c.ClawmachineDatabase.Port,
		c.ClawmachineDatabase.Name,
	)
}

type DiscoveryConfig struct {
	Enabled bool       `mapstructure:"enabled"`
	Etcd    EtcdConfig `mapstructure:"etcd"`
}

type EtcdConfig struct {
	Endpoints []string `mapstructure:"endpoints"`
	Timeout   string   `mapstructure:"timeout"`
}
