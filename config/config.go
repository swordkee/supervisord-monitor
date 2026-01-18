package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var cfgFile string

type SupervisorServer struct {
	URL      string `mapstructure:"url"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type RedmineConfig struct {
	URL        string `mapstructure:"url"`
	AssigneeID string `mapstructure:"assignee_id"`
}

type HTTPAuthConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Config struct {
	SupervisorCols    int                         `mapstructure:"supervisor_cols"`
	Refresh           int                         `mapstructure:"refresh"`
	EnableAlarm       bool                        `mapstructure:"enable_alarm"`
	ShowHost          bool                        `mapstructure:"show_host"`
	Timeout           int                         `mapstructure:"timeout"`
	Port              int                         `mapstructure:"port"`
	SupervisorServers map[string]SupervisorServer `mapstructure:"supervisor_servers"`
	Redmine           RedmineConfig               `mapstructure:"redmine"`
	HTTPAuth          HTTPAuthConfig              `mapstructure:"http_auth"`
}

var Cfg *Config

func LoadConfig(configPath string) error {
	viper.SetConfigType("yaml")

	if configPath != "" {
		absPath, err := filepath.Abs(configPath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}
		viper.SetConfigFile(absPath)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(".", "config"))

		exePath, err := os.Executable()
		if err == nil {
			viper.AddConfigPath(filepath.Dir(exePath))
		}
	}

	viper.SetDefault("supervisor_cols", 2)
	viper.SetDefault("refresh", 10)
	viper.SetDefault("enable_alarm", true)
	viper.SetDefault("show_host", false)
	viper.SetDefault("timeout", 3)
	viper.SetDefault("port", 8080)
	viper.SetDefault("http_auth.username", "")
	viper.SetDefault("http_auth.password", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Config file not found, using defaults: %v", err)
		} else {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	Cfg = &Config{}
	if err := viper.Unmarshal(Cfg); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	return nil
}

func GetServerConfig(serverName string) (*SupervisorServer, error) {
	server, exists := Cfg.SupervisorServers[serverName]
	if !exists {
		return nil, fmt.Errorf("invalid server: %s", serverName)
	}
	return &server, nil
}
