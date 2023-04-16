package internal

import (
	// "fmt"
	"net/http"
	"time"
)

const (
	IdleTimeout  = 60
	MSG_Shutdown = "SHUTDOWN"
)

var (
	_Server *http.Server
)

func init() {
	_Server = &http.Server{
		ReadTimeout:       time.Second * 30,
		WriteTimeout:      time.Minute * 5,
		ReadHeaderTimeout: time.Second * 2,
		MaxHeaderBytes:    1 << 20,
		// Addr:              addr,
		// Handler: engine,
	}
}
