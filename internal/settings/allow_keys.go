package settings

import (
	"fmt"

	"github.com/spf13/viper"
)

type AllowedKeys struct {
	Enable    bool     `mapstructure:"enable"`
	Items     []string `mapstructure:"items"`
	itemsHash map[string]bool
}

func NewAllowedKeys(fp, key string) (cfg *AllowedKeys, err error) {
	vp := viper.New()
	vp.SetConfigType("yaml")

	vp.SetConfigFile(fp)
	if err = vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config %s: %w", fp, err)
	}

	cfg = new(AllowedKeys)
	if err = vp.UnmarshalKey(key, cfg); err != nil {
		return nil, fmt.Errorf("config %s: %w", key, err)
	}

	if !cfg.Enable {
		return cfg, nil
	}

	if len(cfg.Items) == 0 {
		return nil, fmt.Errorf("config %s: items not set", key)
	}

	cfg.itemsHash = make(map[string]bool, len(cfg.Items))
	for i := range cfg.Items {
		cfg.itemsHash[cfg.Items[i]] = true
	}

	return cfg, nil
}

func (cfg *AllowedKeys) Check(item string) bool {
	if !cfg.Enable {
		return true
	}

	return cfg.itemsHash[item]

	/*
		for i := range item.cfg {
			if cfg.Items[i] == item {
				return true
			}
		}

		return false
	*/
}
