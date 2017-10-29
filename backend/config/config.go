package config

import (
	"github.com/BurntSushi/toml"
	"github.com/juju/errors"
)

// Config is used to provide configuration to the server.
type Config struct {
	LogDir   string `toml:"log_dir"`
	Commands struct {
		BasePath      string `toml:"base_path"`
		ConfigCommand string `toml:"config_command"`
	} `toml:"commands"`
	Database struct {
		Driver string `toml:"driver"`
		Spec   string `toml:"spec"`
	} `toml:"database"`
	Server struct {
		Host string `toml:"host"`
		Port uint16 `toml:"port"`
	} `toml:"server"`
	Github struct {
		ClientID     string `toml:"client_id"`
		ClientSecret string `toml:"client_secret"`
	} `toml:"github"`
}

var cfg *Config

func Get() Config {
	if cfg == nil {
		panic("Config is not installed")
	}
	return *cfg
}

func Install(c Config) error {
	switch {
	case c.Commands.BasePath == "":
		return errors.New("Base command path is not configured")
	case c.Commands.ConfigCommand == "":
		return errors.New("Config command is not configured")
	case c.Database.Driver == "":
		return errors.New("Database driver is not configured")
	case c.Database.Spec == "":
		return errors.New("Database spec is not configured")
	case c.Server.Host == "":
		return errors.New("Server host is not configured")
	case c.Server.Port == 0:
		return errors.New("Server port is not configured")
	case c.Github.ClientID == "":
		return errors.New("GitHub client ID is not configured")
	case c.Github.ClientSecret == "":
		return errors.New("GitHub client secret is not configured")
	case c.LogDir == "":
		return errors.New("Log directory is not configured")
	default:
		cfg = &c
		return nil
	}
}

func LoadFile(file string) (Config, error) {
	var c Config
	_, err := toml.DecodeFile(file, &c)
	if err != nil {
		return c, errors.Annotate(err, "Failed to parse config file")
	}
	return c, Install(c)
}
