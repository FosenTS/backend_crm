package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type AppConfig struct {
	Server struct {
		Host         string `json:"host"`
		Port         int    `json:"port"`
		ReadTimeout  string `json:"read_timeout"`
		WriteTimeout string `json:"write_timeout"`
	} `json:"server"`

	TLS struct {
		CertFilePath string `json:"cert_file_path"`
		CertKeyPath  string `json:"cert_key_path"`
	} `json:"tls"`

	JWT struct {
		AccessSecret  string `json:"access_secret"`
		RefreshSecret string `json:"refresh_secret"`
		AccessTTL     string `json:"access_ttl"`
		RefreshTTL    string `json:"refresh_ttl"`
	} `json:"jwt"`

	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"db_name"`
		SSLMode  string `json:"ssl_mode"`
	} `json:"database"`

	HTML struct {
		BasePath string `json:"base_path"`
		Files    struct {
			Index    string `json:"index"`
			Login    string `json:"login"`
			Register string `json:"register"`
			Orders   string `json:"orders"`
		} `json:"files"`
	} `json:"html"`

	// Parsed durations
	parsedReadTimeout  time.Duration
	parsedWriteTimeout time.Duration
	parsedAccessTTL    time.Duration
	parsedRefreshTTL   time.Duration
}

func NewConfig() (*AppConfig, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.json"
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config AppConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	// Parse durations
	if config.Server.ReadTimeout == "" {
		config.Server.ReadTimeout = "5s"
	}
	if config.Server.WriteTimeout == "" {
		config.Server.WriteTimeout = "5s"
	}
	if config.JWT.AccessTTL == "" {
		config.JWT.AccessTTL = "15m"
	}
	if config.JWT.RefreshTTL == "" {
		config.JWT.RefreshTTL = "720h"
	}

	var parseErr error
	config.parsedReadTimeout, parseErr = time.ParseDuration(config.Server.ReadTimeout)
	if parseErr != nil {
		return nil, parseErr
	}

	config.parsedWriteTimeout, parseErr = time.ParseDuration(config.Server.WriteTimeout)
	if parseErr != nil {
		return nil, parseErr
	}

	config.parsedAccessTTL, parseErr = time.ParseDuration(config.JWT.AccessTTL)
	if parseErr != nil {
		return nil, parseErr
	}

	config.parsedRefreshTTL, parseErr = time.ParseDuration(config.JWT.RefreshTTL)
	if parseErr != nil {
		return nil, parseErr
	}

	// Set default values if not provided
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}

	if config.Database.Port == 0 {
		config.Database.Port = 5432
	}
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}

	// Ensure HTML file paths are absolute
	if config.HTML.BasePath != "" {
		config.HTML.Files.Index = filepath.Join(config.HTML.BasePath, config.HTML.Files.Index)
		config.HTML.Files.Login = filepath.Join(config.HTML.BasePath, config.HTML.Files.Login)
		config.HTML.Files.Register = filepath.Join(config.HTML.BasePath, config.HTML.Files.Register)
		config.HTML.Files.Orders = filepath.Join(config.HTML.BasePath, config.HTML.Files.Orders)
	}

	return &config, nil
}

func (c *AppConfig) GetDSN() string {
	return "postgres://" + c.Database.User + ":" + c.Database.Password + "@" + c.Database.Host + ":" + strconv.Itoa(c.Database.Port) + "/" + c.Database.DBName + "?sslmode=" + c.Database.SSLMode
}

func (c *AppConfig) GetServerAddr() string {
	return c.Server.Host + ":" + strconv.Itoa(c.Server.Port)
}

// GetReadTimeout returns the parsed read timeout duration
func (c *AppConfig) GetReadTimeout() time.Duration {
	return c.parsedReadTimeout
}

// GetWriteTimeout returns the parsed write timeout duration
func (c *AppConfig) GetWriteTimeout() time.Duration {
	return c.parsedWriteTimeout
}

// GetAccessTTL returns the parsed access token TTL duration
func (c *AppConfig) GetAccessTTL() time.Duration {
	return c.parsedAccessTTL
}

// GetRefreshTTL returns the parsed refresh token TTL duration
func (c *AppConfig) GetRefreshTTL() time.Duration {
	return c.parsedRefreshTTL
}
