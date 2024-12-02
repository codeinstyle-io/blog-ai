package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
)

type Config struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}
	DB struct {
		Path     string `mapstructure:"path"`
		LogLevel string `mapstructure:"log_level"`
	}
	Site struct {
		ChromaStyle  string `mapstructure:"chroma_style"`
		Timezone     string `mapstructure:"timezone"`
		Title        string `mapstructure:"title"`
		Subtitle     string `mapstructure:"subtitle"`
		Theme        string `mapstructure:"theme"`
		PostsPerPage int    `mapstructure:"posts_per_page"`
	}
}

func InitConfig() (*Config, error) {
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("db.path", "blog.db")
	viper.SetDefault("db.log_level", "warn")
	viper.SetDefault("site.chroma_style", "paraiso-dark")
	viper.SetDefault("site.timezone", "UTC")
	viper.SetDefault("site.title", "Captain")
	viper.SetDefault("site.subtitle", "A simple blog engine")
	viper.SetDefault("site.theme", "")
	viper.SetDefault("site.posts_per_page", 3)

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
		fmt.Printf("Loaded config from %s\n", configFile)
	} else {
		fmt.Println("Using default config")
	}

	return &cfg, nil
}

// GetGormLogLevel returns the gorm logger level based on the config
func (c *Config) GetGormLogLevel() logger.LogLevel {
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

// GetLocation returns the configured timezone location
func (c *Config) GetLocation() *time.Location {
	loc, err := time.LoadLocation(c.Site.Timezone)
	if err != nil {
		return time.UTC
	}
	return loc
}
