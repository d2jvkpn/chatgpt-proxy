package settings

import (
	// "fmt"

	"github.com/d2jvkpn/x-ai/pkg/lang_chain"
)

type LangChainAgent struct {
	*lang_chain.LangChain
}

func NewLangChainAgent(key, path string) (lca *LangChainAgent, err error) {
	lca = new(LangChainAgent)

	if lca.LangChain, err = lang_chain.NewLangChain(key, path); err != nil {
		return nil, err
	}

	return lca, nil
}
