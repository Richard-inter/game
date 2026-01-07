package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Default configuration constants
const (
	defaultServerPort    = 8080
	defaultServerTimeout = 30
	defaultDatabasePort  = 3306
	defaultRedisPort     = 6379
	defaultGRPCPort      = 9090
	defaultWebSocketPort = 8081
	defaultTCPPort       = 8082
	defaultBufferSize    = 1024
	defaultJWTExpiration = 86400 // 24 hours
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

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Set environment variable prefix
	viper.SetEnvPrefix("GAME")
	viper.AutomaticEnv()

	// Set default values
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found, using defaults and environment variables")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", defaultServerPort)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", defaultServerTimeout)
	viper.SetDefault("server.write_timeout", defaultServerTimeout)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", defaultDatabasePort)
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.name", "game")
	viper.SetDefault("database.charset", "utf8mb4")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", defaultRedisPort)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// gRPC defaults
	viper.SetDefault("grpc.host", "0.0.0.0")
	viper.SetDefault("grpc.port", defaultGRPCPort)

	// WebSocket defaults
	viper.SetDefault("websocket.host", "0.0.0.0")
	viper.SetDefault("websocket.port", defaultWebSocketPort)
	viper.SetDefault("websocket.path", "/ws")
	viper.SetDefault("websocket.read_buffer_size", defaultBufferSize)
	viper.SetDefault("websocket.write_buffer_size", defaultBufferSize)

	// TCP defaults
	viper.SetDefault("tcp.host", "0.0.0.0")
	viper.SetDefault("tcp.port", defaultTCPPort)
	viper.SetDefault("tcp.keep_alive", true)
	viper.SetDefault("tcp.read_timeout", defaultServerTimeout)
	viper.SetDefault("tcp.write_timeout", defaultServerTimeout)

	// JWT defaults
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expiration_time", defaultJWTExpiration) // 24 hours

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")

	// Tracing defaults
	viper.SetDefault("tracing.enabled", false)
	viper.SetDefault("tracing.service_name", "game-server")
	viper.SetDefault("tracing.jaeger_url", "http://localhost:14268/api/traces")
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
