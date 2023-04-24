package internal

import (
	"fmt"

	"github.com/spf13/viper"
)

type TlsConfig struct {
	Enable bool   `mapstructure:"enable"`
	Cert   string `mapstructure:"cert"`
	Key    string `mapstructure:"key"`
}

func NewTlsConfig(fp, key string) (config *TlsConfig, err error) {
	vp := viper.New()
	vp.SetConfigType("yaml")

	vp.SetConfigFile(fp)
	if err = vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ReadInConfig(): %q, %w", fp, err)
	}

	config = new(TlsConfig)
	if err = vp.UnmarshalKey(key, config); err != nil {
		return nil, err
	}

	return config, nil
}
