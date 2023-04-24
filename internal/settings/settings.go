package settings

import (
	// "fmt"
	"strings"

	"github.com/d2jvkpn/chatgpt-proxy/internal/biz"

	"github.com/d2jvkpn/go-web/pkg/wrap"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	_Project *viper.Viper

	// GPTCli       *chatgpt.Client
	GPTCli2 *openai.Client
	LCA     *biz.LangChainAgent

	AllowIps     *AllowedKeys
	AllowApiKeys *AllowedKeys

	Logger      *wrap.Logger
	TransLogger *zap.Logger // transaction
	ReqLogger   *zap.Logger // request
	AppLogger   *zap.Logger // application
	// DebugLogger    *zap.Logger // debug
)

func init() {
	// StartTs = time.Now().Unix()
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
