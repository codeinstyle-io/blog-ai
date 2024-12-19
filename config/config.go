package config

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
)

var (
	// Available timezones
	timezones = []string{
		"UTC",
		"Europe/London",
		"Europe/Paris",
		"Europe/Berlin",
		"Europe/Madrid",
		"Europe/Rome",
		"Europe/Amsterdam",
		"Europe/Brussels",
		"Europe/Vienna",
		"Europe/Stockholm",
		"Europe/Copenhagen",
		"Europe/Oslo",
		"Europe/Warsaw",
		"Europe/Moscow",
		"Europe/Istanbul",
		"America/New_York",
		"America/Chicago",
		"America/Denver",
		"America/Los_Angeles",
		"America/Toronto",
		"America/Vancouver",
		"America/Mexico_City",
		"America/Sao_Paulo",
		"America/Buenos_Aires",
		"Asia/Dubai",
		"Asia/Tokyo",
		"Asia/Shanghai",
		"Asia/Hong_Kong",
		"Asia/Singapore",
		"Asia/Seoul",
		"Asia/Kolkata",
		"Asia/Bangkok",
		"Asia/Jakarta",
		"Australia/Sydney",
		"Australia/Melbourne",
		"Australia/Brisbane",
		"Australia/Perth",
		"Pacific/Auckland",
	}

	// Available Chroma styles
	chromaStyles = []string{
		"abap",
		"algol",
		"arduino",
		"autumn",
		"borland",
		"bw",
		"colorful",
		"dracula",
		"emacs",
		"friendly",
		"fruity",
		"github",
		"igor",
		"lovelace",
		"manni",
		"monokai",
		"monokailight",
		"murphy",
		"native",
		"paraiso-dark",
		"paraiso-light",
		"pastie",
		"perldoc",
		"pygments",
		"rainbow_dash",
		"rrt",
		"solarized-dark",
		"solarized-dark256",
		"solarized-light",
		"swapoff",
		"tango",
		"trac",
		"vim",
		"vs",
		"xcode",
	}
)

type Config struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}
	Site struct {
		SecureCookie bool   `mapstructure:"secure_cookie"`
		Domain       string `mapstructure:"domain"`
		Theme        string `mapstructure:"theme"`
	} `mapstructure:"site"`
	DB struct {
		Path     string `mapstructure:"path"`
		LogLevel string `mapstructure:"log_level"`
	}
	Storage struct {
		Provider string `mapstructure:"provider"` // "local" or "s3"
		S3       struct {
			Bucket    string `mapstructure:"bucket"`
			Region    string `mapstructure:"region"`
			Endpoint  string `mapstructure:"endpoint"`
			AccessKey string `mapstructure:"access_key"`
			SecretKey string `mapstructure:"secret_key"`
		}
		LocalPath string `mapstructure:"local_path"` // Path for local storage
	}
	Debug bool `mapstructure:"debug"`
}

func InitConfig() (*Config, error) {
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("site.secure_cookie", false)
	viper.SetDefault("site.domain", "")
	viper.SetDefault("site.theme", "")
	viper.SetDefault("db.path", "blog.db")
	viper.SetDefault("db.log_level", "warn")
	viper.SetDefault("storage.provider", "local")
	viper.SetDefault("storage.local_path", "./storage")

	// S3
	viper.SetDefault("storage.s3.bucket", "")
	viper.SetDefault("storage.s3.region", "")
	viper.SetDefault("storage.s3.endpoint", "")
	viper.SetDefault("storage.s3.access_key", "")
	viper.SetDefault("storage.s3.secret_key", "")

	// Debug
	viper.SetDefault("debug", false)

	// Read config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/captain/")

	// Enable environment variables
	viper.SetEnvPrefix("CAPTAIN")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read config file if exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	configFile := viper.ConfigFileUsed()

	if configFile != "" {
		log.Infof("Loaded config from %s\n", configFile)
	} else {
		log.Debug("Using default config")
	}

	return &cfg, nil
}

// GetGormLogLevel returns the gorm logger level based on the config
func (c *Config) GetGormLogLevel() logger.LogLevel {
	if c.Debug {
		return logger.Info
	}

	switch strings.ToLower(c.DB.LogLevel) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
}

// GetTimezones returns the list of available timezones
func (c *Config) GetTimezones() []string {
	return timezones
}

// GetChromaStyles returns the list of available syntax highlighting themes
func (c *Config) GetChromaStyles() []string {
	return chromaStyles
}

// ValidateS3Config validates S3 configuration if S3 provider is selected
func (c *Config) ValidateS3Config() error {
	if c.Storage.Provider != "s3" {
		return nil
	}

	var missingFields []string

	if c.Storage.S3.Bucket == "" {
		missingFields = append(missingFields, "storage.s3.bucket")
	}
	if c.Storage.S3.Region == "" {
		missingFields = append(missingFields, "storage.s3.region")
	}
	if c.Storage.S3.AccessKey == "" {
		missingFields = append(missingFields, "storage.s3.access_key")
	}
	if c.Storage.S3.SecretKey == "" {
		missingFields = append(missingFields, "storage.s3.secret_key")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required S3 configuration fields: %v", missingFields)
	}

	return nil
}
