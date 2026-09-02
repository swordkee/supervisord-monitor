package config

import (
	"crypto/subtle"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var cfgFile string

type SupervisorServer struct {
	Name     string `mapstructure:"name"`
	URL      string `mapstructure:"url"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type RedmineConfig struct {
	URL        string `mapstructure:"url"`
	AssigneeID string `mapstructure:"assignee_id"`
}

// MonitorGroup is a named collection of supervisor servers that can be
// granted to web users as a unit.
type MonitorGroup struct {
	Name    string   `mapstructure:"name"`
	Servers []string `mapstructure:"servers"`
}

// AuthUser is a web interface login account. Access is derived from the
// listed monitor groups and/or explicit server names. A "*" entry in either
// field grants access to every server.
type AuthUser struct {
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
	Groups   []string `mapstructure:"groups"`
	Servers  []string `mapstructure:"servers"`
}

type HTTPAuthConfig struct {
	Username string     `mapstructure:"username"`
	Password string     `mapstructure:"password"`
	Users    []AuthUser `mapstructure:"users"`
}

// AuthPrincipal represents an authenticated web user and the set of servers
// it is allowed to see and control. A nil AllowedServers map means the user
// may access every configured server.
type AuthPrincipal struct {
	Username       string
	AllowedServers map[string]bool
}

type Config struct {
	SupervisorCols    int                `mapstructure:"supervisor_cols"`
	Refresh           int                `mapstructure:"refresh"`
	EnableAlarm       bool               `mapstructure:"enable_alarm"`
	ShowHost          bool               `mapstructure:"show_host"`
	Timeout           int                `mapstructure:"timeout"`
	Port              int                `mapstructure:"port"`
	SupervisorServers []SupervisorServer `mapstructure:"supervisor_servers"`
	Redmine           RedmineConfig      `mapstructure:"redmine"`
	MonitorGroups     []MonitorGroup     `mapstructure:"monitor_groups"`
	HTTPAuth          HTTPAuthConfig     `mapstructure:"http_auth"`
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

// AuthEnabled reports whether web interface authentication is configured.
func (c *Config) AuthEnabled() bool {
	return (c.HTTPAuth.Username != "" && c.HTTPAuth.Password != "") || len(c.HTTPAuth.Users) > 0
}

// Authenticate validates web interface credentials and returns the principal
// describing which servers the account may access. It returns nil when the
// credentials do not match any configured account.
func (c *Config) Authenticate(username, password string) *AuthPrincipal {
	groupServers := make(map[string][]string, len(c.MonitorGroups))
	for i := range c.MonitorGroups {
		group := &c.MonitorGroups[i]
		groupServers[group.Name] = group.Servers
	}

	// Legacy single-account configuration: grants access to every server.
	if c.HTTPAuth.Username != "" && c.HTTPAuth.Password != "" {
		userOK := subtle.ConstantTimeCompare([]byte(username), []byte(c.HTTPAuth.Username)) == 1
		passOK := subtle.ConstantTimeCompare([]byte(password), []byte(c.HTTPAuth.Password)) == 1
		if userOK && passOK {
			return &AuthPrincipal{Username: c.HTTPAuth.Username, AllowedServers: nil}
		}
	}

	for i := range c.HTTPAuth.Users {
		user := &c.HTTPAuth.Users[i]
		if user.Username == "" || user.Password == "" {
			continue
		}

		userOK := subtle.ConstantTimeCompare([]byte(username), []byte(user.Username)) == 1
		passOK := subtle.ConstantTimeCompare([]byte(password), []byte(user.Password)) == 1
		if !userOK || !passOK {
			continue
		}

		principal := &AuthPrincipal{Username: user.Username}
		allowed := make(map[string]bool)
		grantsAll := false

		for _, groupName := range user.Groups {
			if groupName == "*" {
				grantsAll = true
				break
			}
			for _, srv := range groupServers[groupName] {
				allowed[srv] = true
			}
		}

		if !grantsAll {
			for _, srv := range user.Servers {
				if srv == "*" {
					grantsAll = true
					break
				}
				allowed[srv] = true
			}
		}

		if grantsAll {
			principal.AllowedServers = nil
		} else {
			principal.AllowedServers = allowed
		}
		return principal
	}

	return nil
}

func GetServerConfig(serverName string) (*SupervisorServer, error) {
	for i := range Cfg.SupervisorServers {
		if Cfg.SupervisorServers[i].Name == serverName {
			return &Cfg.SupervisorServers[i], nil
		}
	}
	return nil, fmt.Errorf("invalid server: %s", serverName)
}

func GetOrderedServerNames() []string {
	names := make([]string, len(Cfg.SupervisorServers))
	for i, server := range Cfg.SupervisorServers {
		names[i] = server.Name
	}
	return names
}

func GetServerByName(name string) *SupervisorServer {
	for i := range Cfg.SupervisorServers {
		if Cfg.SupervisorServers[i].Name == name {
			return &Cfg.SupervisorServers[i]
		}
	}
	return nil
}
