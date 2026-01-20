package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Configuration constants for validation
const (
	minServerPort    = 1024
	maxServerPort    = 65535
	minDatabasePort  = 1
	maxDatabasePort  = 65535
	minRedisPort     = 1
	maxRedisPort     = 65535
	minGRPCPort      = 1
	maxGRPCPort      = 65535
	minWebSocketPort = 1
	maxWebSocketPort = 65535
	minTCPPort       = 1
	maxTCPPort       = 65535
	minBufferSize    = 512
	maxBufferSize    = 65536
	minJWTExpiration = 300        // 5 minutes
	maxJWTExpiration = 86400 * 30 // 30 days
)

type RPCServiceConfig struct {
	Enabled            bool   `mapstructure:"enabled"`
	ServicePackage     string `mapstructure:"service_package"`
	ServiceName        string `mapstructure:"service_name"`
	ImportPath         string `mapstructure:"import_path"`
	ImplementationPath string `mapstructure:"implementation_path"`
}

type RPCServicesConfig struct {
	ClawMachine RPCServiceConfig `mapstructure:"claw_machine"`
	Player      RPCServiceConfig `mapstructure:"player"`
}

type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Redis       RedisConfig       `mapstructure:"redis"`
	GRPC        GRPCConfig        `mapstructure:"grpc"`
	RPCServices RPCServicesConfig `mapstructure:"rpc_services"`
	WebSocket   WebSocketConfig   `mapstructure:"websocket"`
	TCP         TCPConfig         `mapstructure:"tcp"`
	JWT         JWTConfig         `mapstructure:"jwt"`
	Logging     LoggingConfig     `mapstructure:"logging"`
	Tracing     TracingConfig     `mapstructure:"tracing"`
}

type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Charset  string `mapstructure:"charset"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type GRPCConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type WebSocketConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Path            string `mapstructure:"path"`
	ReadBufferSize  int    `mapstructure:"read_buffer_size"`
	WriteBufferSize int    `mapstructure:"write_buffer_size"`
}

type TCPConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	KeepAlive    bool   `mapstructure:"keep_alive"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type JWTConfig struct {
	Secret         string `mapstructure:"secret"`
	ExpirationTime int    `mapstructure:"expiration_time"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type TracingConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	ServiceName string `mapstructure:"service_name"`
	JaegerURL   string `mapstructure:"jaeger_url"`
}

type StreamConsumerConfig struct {
	StreamKey     string `mapstructure:"stream_key"`
	ConsumerGroup string `mapstructure:"consumer_group"`
	ConsumerName  string `mapstructure:"consumer_name"`
	BatchSize     int    `mapstructure:"batch_size"`
	BlockTimeout  string `mapstructure:"block_timeout"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Set environment variable prefix
	viper.SetEnvPrefix("GAME")
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found and no defaults provided")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &config, nil
}

func validateConfig(config *Config) error {
	// Validate server configuration
	if config.Server.Port < minServerPort || config.Server.Port > maxServerPort {
		return fmt.Errorf("server port must be between %d and %d", minServerPort, maxServerPort)
	}

	// Validate database configuration
	if config.Database.Port < minDatabasePort || config.Database.Port > maxDatabasePort {
		return fmt.Errorf("database port must be between %d and %d", minDatabasePort, maxDatabasePort)
	}

	// Validate Redis configuration
	if config.Redis.Port < minRedisPort || config.Redis.Port > maxRedisPort {
		return fmt.Errorf("redis port must be between %d and %d", minRedisPort, maxRedisPort)
	}

	// Validate gRPC configuration
	if config.GRPC.Port < minGRPCPort || config.GRPC.Port > maxGRPCPort {
		return fmt.Errorf("grpc port must be between %d and %d", minGRPCPort, maxGRPCPort)
	}

	// Validate WebSocket configuration
	if config.WebSocket.Port < minWebSocketPort || config.WebSocket.Port > maxWebSocketPort {
		return fmt.Errorf("websocket port must be between %d and %d", minWebSocketPort, maxWebSocketPort)
	}

	if config.WebSocket.ReadBufferSize < minBufferSize || config.WebSocket.ReadBufferSize > maxBufferSize {
		return fmt.Errorf("websocket read buffer size must be between %d and %d", minBufferSize, maxBufferSize)
	}

	if config.WebSocket.WriteBufferSize < minBufferSize || config.WebSocket.WriteBufferSize > maxBufferSize {
		return fmt.Errorf("websocket write buffer size must be between %d and %d", minBufferSize, maxBufferSize)
	}

	// Validate TCP configuration
	if config.TCP.Port < minTCPPort || config.TCP.Port > maxTCPPort {
		return fmt.Errorf("tcp port must be between %d and %d", minTCPPort, maxTCPPort)
	}

	// Validate JWT configuration
	if config.JWT.ExpirationTime < minJWTExpiration || config.JWT.ExpirationTime > maxJWTExpiration {
		return fmt.Errorf("jwt expiration time must be between %d and %d seconds", minJWTExpiration, maxJWTExpiration)
	}

	if config.JWT.Secret == "" {
		return fmt.Errorf("jwt secret cannot be empty")
	}

	return nil
}

// GetDSN returns database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.Charset,
	)
}

// GetRedisAddr returns Redis connection address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// GetServerAddr returns server address
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetGRPCAddr returns gRPC server address
func (c *Config) GetGRPCAddr() string {
	return fmt.Sprintf("%s:%d", c.GRPC.Host, c.GRPC.Port)
}

// GetWebSocketAddr returns WebSocket server address
func (c *Config) GetWebSocketAddr() string {
	return fmt.Sprintf("%s:%d", c.WebSocket.Host, c.WebSocket.Port)
}

// GetTCPAddr returns TCP server address
func (c *Config) GetTCPAddr() string {
	return fmt.Sprintf("%s:%d", c.TCP.Host, c.TCP.Port)
}
