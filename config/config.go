package config

import "time"

type Config struct {
	ChromaStyle string
	Timezone    *time.Location
}

func NewDefaultConfig() *Config {
	location, _ := time.LoadLocation("UTC")
	return &Config{
		ChromaStyle: "paraiso-dark",
		Timezone:    location,
	}
}
