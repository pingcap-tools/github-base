package config

import (
	"github.com/BurntSushi/toml"
	"github.com/juju/errors"
)

// Config contains configuration options.
type Config struct {
	Host         string      `toml:"host"`
	Port         int         `toml:"port"`
	GithubToken  string      `toml:"github"`
	GithubSecret string      
	Database     *Database   `toml:"database"`
}

// Database defines db configuration
type Database struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

var globalConf = Config{
	Host: "0.0.0.0",
	Port: 30000,
	Database: &Database{
		Host: "127.0.0.1",
		Port: 3306,
		Username: "root",
		Password: "",
		Database: "base",
	},
}

// GetGlobalConfig returns the global configuration for this server.
func GetGlobalConfig() *Config {
	return &globalConf
}

// Load loads config options from a toml file_logger.
func (c *Config)Load(confFile string) error {
	_, err := toml.DecodeFile(confFile, c)
	return errors.Trace(err)
}

// Init do some prepare works
func (c *Config)Init() error {
	return nil
}
