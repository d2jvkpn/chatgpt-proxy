package settings

import (
	"fmt"
	"strings"

	"github.com/d2jvkpn/x-ai/pkg/chatgpt"

	"github.com/spf13/viper"
)

var (
	GPTCli       *chatgpt.Client
	_Project     *viper.Viper
	AllowIps     *AllowedKeys
	AllowApiKeys *AllowedKeys
	Tls          *TlsConfig
)

type TlsConfig struct {
	Enable bool   `mapstructure:"enable"`
	Crt    string `mapstructure:"crt"`
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

func SetProject(str string) (err error) {
	_Project = viper.New()
	_Project.SetConfigType("yaml")
	return _Project.ReadConfig(strings.NewReader(str))
}

func Version() string {
	return _Project.GetString("version")
}

func Config() string {
	return _Project.GetString("config")
}
