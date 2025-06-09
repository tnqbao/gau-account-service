package config

type Config struct {
	EnvConfig *EnvConfig `json:"env_config"`
}

func NewConfig() *Config {
	EnvConfig := LoadEnvConfig()
	return &Config{
		EnvConfig: EnvConfig,
	}
}
