package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/d2jvkpn/chatgpt-proxy/internal"
	"github.com/d2jvkpn/chatgpt-proxy/internal/settings"

	"github.com/d2jvkpn/go-web/pkg/misc"
)

var (
	//go:embed project.yaml
	_Project    string
	_NotWindows bool
)

func init() {
	_NotWindows = runtime.GOOS != "windows"
}

func main() {
	var (
		release bool
		config  string
		addr    string
		err     error
		meta    map[string]any
		errch   chan error
		quit    chan os.Signal
	)

	if err = settings.SetProject(_Project); err != nil {
		log.Fatalln(err)
	}

	meta = misc.BuildInfo()
	meta["project"] = settings.GetProject()
	meta["version"] = settings.GetVersion()

	flag.StringVar(&addr, "addr", "0.0.0.0:3021", "http server listening address")
	flag.StringVar(&config, "config", "configs/local.yaml", "config file path")
	flag.BoolVar(&release, "release", false, "run in release mode")

	flag.Usage = func() {
		output := flag.CommandLine.Output()

		fmt.Fprintf(output, "Usage:\n")
		flag.PrintDefaults()
		fmt.Fprintf(output, "\nConfig:\n```yaml\n%s```\n", settings.GetConfig())
	}

	flag.Parse()

	if err = internal.Load(config, release); err != nil {
		log.Fatalln(err)
	}

	meta["release"] = release
	if errch, err = internal.Serve(addr, meta); err != nil {
		log.Fatalln(err)
	}
	protocol := "http"
	if internal.TlsEnabled() {
		protocol = "https"
	}

	fmt.Printf(">>> HTTP server is listening on %s://%s\n", protocol, addr)

	quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR2)

	select {
	case err = <-errch:
		// break
	case <-quit: // sig := <-quit:
		// if sig == syscall.SIGUSR2 {}
		fmt.Println("")
		errch <- fmt.Errorf(internal.MSG_Shutdown)
		log.Printf("<<< Exit\n")
		<-errch
	}

	if err != nil {
		log.Fatalln(err)
	}
}
