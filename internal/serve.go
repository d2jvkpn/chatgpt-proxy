package internal

import (
	// "fmt"
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/d2jvkpn/chatgpt-proxy/internal/settings"
)

func Serve(addr string, release bool) (errch chan error, err error) {
	var (
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
			shutdown()
		} else {
			errch <- err
		}
	}()

	go func() {
		if err := <-errch; err.Error() == "SHUTDOWN" {
			shutdown()
		}
	}()

	return errch, nil
}

func shutdown() {
	var err error

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	if err = _Server.Shutdown(ctx); err != nil {
		// log.Error(fmt.Sprintf("server shutdown: %v", err))
	}
	cancel()
}
