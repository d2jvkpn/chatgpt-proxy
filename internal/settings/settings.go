package settings

import (
	"fmt"
	"strings"

	"github.com/d2jvkpn/go-web/pkg/wrap"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	_Project *viper.Viper

	// GPTCli       *chatgpt.Client
	GPTCli2 *openai.Client
	LCA     *LangChainAgent

	AllowIps     *AllowedKeys
	AllowApiKeys *AllowedKeys
	Tls          *TlsConfig

	Logger      *wrap.Logger
	TransLogger *zap.Logger // transaction
	ReqLogger   *zap.Logger // request
	AppLogger   *zap.Logger // application
	// DebugLogger    *zap.Logger // debug
)

type TlsConfig struct {
	Enable bool   `mapstructure:"enable"`
	Crt    string `mapstructure:"crt"`
	Key    string `mapstructure:"key"`
}

func init() {
	// StartTs = time.Now().Unix()
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

func GetProject() string {
	return _Project.GetString("project")
}

func GetVersion() string {
	return _Project.GetString("version")
}

func GetConfig() string {
	return _Project.GetString("config")
}

func SetupLoggers() {
	TransLogger = Logger.Named("transaction")
	ReqLogger = Logger.Named("request")
	AppLogger = Logger.Named("application")
}
