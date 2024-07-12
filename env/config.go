package env

import (
	"github.com/spf13/viper"
	"spy-cat-api/services"
	"time"
)

const (
	defaultHttpPort      = "8080"
	defaultHttpRWTimeout = 10 * time.Second

	defaultConfigPath = "./config.yml"
)

type Config struct {
	Http services.HttpServerConfig
	DB   services.DBConfig
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = defaultConfigPath
	}
	populateDefaults()

	if err := parseConfigFile(path); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func populateDefaults() {
	viper.SetDefault("http.host", "localhost:"+defaultHttpPort)
	viper.SetDefault("http.readTimeout", defaultHttpRWTimeout)
	viper.SetDefault("http.writeTimeout", defaultHttpRWTimeout)
}

func parseConfigFile(filePath string) error {
	viper.SetConfigFile(filePath)

	return viper.ReadInConfig()
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.Http); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("db", &cfg.DB); err != nil {
		return err
	}

	return nil
}
