package config

type Config struct {
	ChromaStyle string
}

func NewDefaultConfig() *Config {
	return &Config{
		ChromaStyle: "paraiso-dark",
	}
}
