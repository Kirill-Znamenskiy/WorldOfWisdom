package config

import (
	"context"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	"github.com/Kirill-Znamenskiy/kzfunc/fn"
)

type Ctx = context.Context

const (
	ConfigFileEnvKey      = "CONFIG_FILE"
	ConfigFileDefaultPath = "./config.yaml"

	ENVDEV   = "DEV"
	ENVPROD  = "PROD"
	ENVLOCAL = "LOCAL"
)

func Init(ctx Ctx) (ret *Config, err error) {
	viper.New()

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "___"))

	cfpaths, err := DetectConfigFilePaths()
	if err != nil {
		return nil, err
	}

	for _, cfpath := range cfpaths {
		cmp, err := LoadConfigMapFromFile(cfpath)
		if err != nil {
			return nil, err
		}

		err = viper.MergeConfigMap(cmp)
		if err != nil {
			return nil, err
		}
	}

	allViperSettings := viper.AllSettings()
	_ = allViperSettings

	ret = new(Config)
	err = viper.Unmarshal(ret)
	if err != nil {
		return nil, err
	}

	err = validator.New().StructCtx(ctx, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func DetectConfigFilePaths() (ret []string, err error) {
	path := os.Getenv(ConfigFileEnvKey)

	if path == "" && len(os.Args) > 1 {
		path = os.Args[1]
	}

	if path == "" {
		path = ConfigFileDefaultPath
	}

	paths := strings.Split(path, "+")

	ret = fn.MapFilter(paths, func(path string) *string {
		path = strings.TrimSpace(path)
		if path == "" {
			return nil
		}
		return &path
	})

	return ret, nil
}

func LoadConfigMapFromFile(cfpath string) (ret map[string]any, err error) {
	tmp := viper.New()
	tmp.SetConfigFile(cfpath)

	err = tmp.MergeInConfig()
	if err != nil {
		return nil, err
	}

	ret = tmp.AllSettings()

	return ret, nil
}
