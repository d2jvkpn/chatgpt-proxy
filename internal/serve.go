package internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/d2jvkpn/chatgpt-proxy/internal/settings"

	"go.uber.org/zap"
)

func Serve(addr string, meta map[string]any) (errch chan error, err error) {
	var (
		once     *sync.Once
		listener net.Listener
		cert     tls.Certificate
	)

	if listener, err = net.Listen("tcp", addr); err != nil {
		return nil, err
	}

	if settings.Tls.Enable {
		cert, err = tls.LoadX509KeyPair(settings.Tls.Crt, settings.Tls.Key)
		if err != nil {
			return nil, err
		}

		_Server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	meta["allowIps"] = settings.AllowIps.Enable
	meta["apiKeys"] = settings.AllowApiKeys.Enable
	meta["tls"] = settings.Tls.Enable

	settings.AppLogger.Info("the server is starting", zap.Any("meta", meta))

	once = new(sync.Once)

	shutdown := func() {
		var err error

		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		if err = _Server.Shutdown(ctx); err != nil {
			settings.AppLogger.Error(fmt.Sprintf("shutdown the server : %v", err))
		} else {
			settings.AppLogger.Warn("the server is shutting down")
		}
		cancel()
	}

	errch = make(chan error, 2)
	go func() {
		// err := _Server.ServeTLS(listener, "configs/localhost.csr", "configs/localhost.key")
		var err error

		if _Server.TLSConfig == nil {
			err = _Server.Serve(listener)
		} else {
			err = _Server.ServeTLS(listener, "", "")
		}

		if err != http.ErrServerClosed {
			once.Do(onexit)
			errch <- err
		}
	}()

	go func() {
		var err = <-errch
		if err.Error() == MSG_Shutdown {
			shutdown()
			once.Do(onexit)
			errch <- nil
		}
	}()

	return errch, nil
}

func onexit() {
	settings.Logger.Down()
}
