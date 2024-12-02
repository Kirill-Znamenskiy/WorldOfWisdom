package config

import "strings"

type Config struct {
	ENV string `mapstructure:"env" validate:"required" desc:"(env ENV) environment: LOCAL, DEV or PROD"`

	//nolint:lll // can't separate tags
	LogLevel string `mapstructure:"log_level" validate:"required" desc:"(env LOG_LEVEL) log level: DEBUG,INFO(default),WARN,ERROR,CRITICAL"`

	BuildGitShowVersion string `mapstructure:"-" desc:"build app version in form of git show command"`

	Server struct {
		Address string `mapstructure:"address" validate:"required" desc:"server address"`
	} `mapstructure:"server" validate:"required"`
}

func (c *Config) IsDEV() bool   { return strings.ToUpper(c.ENV) == ENVDEV }
func (c *Config) IsPROD() bool  { return strings.ToUpper(c.ENV) == ENVPROD }
func (c *Config) IsLOCAL() bool { return strings.ToUpper(c.ENV) == ENVLOCAL }