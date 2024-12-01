package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
)

type Config struct {
	Server struct {
		Host string
		Port int
	}
	DB struct {
		Path     string
		LogLevel string
	}
	Site struct {
		ChromaStyle string
		Timezone    string
		Title       string
		Subtitle    string
		Theme       string
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
		// Only return error if it's not a "config file not found" error
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	err := viper.Unmarshal(&config)
	return &config, err
}

func (c *Config) GetGormLogLevel() logger.LogLevel {
	switch strings.ToLower(c.DB.LogLevel) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
}

func (c *Config) GetLocation() *time.Location {
	loc, err := time.LoadLocation(c.Site.Timezone)
	if err != nil {
		loc, _ = time.LoadLocation("UTC")
	}
	return loc
}
