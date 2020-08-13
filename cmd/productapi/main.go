package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"

	"github.com/mbellgb/productsvc"
)

func main() {
	var (
		httpHost = flag.String("http.host", "127.0.0.1", "Host to bind HTTP server to")
		httpPort = flag.String("http.port", "8080", "Port to serve HTTP server on")
	)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	logger = log.With(logger, "application", "productsvc")

	var svc productsvc.Service = productsvc.NewInMemoryService()

	var handler http.Handler = productsvc.MakeHTTPHandler(svc, log.NewJSONLogger(os.Stdout))

	errs := make(chan error)

	// Goroutine to listen for termination signals from OS.
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// Main server goroutine.
	go func() {
		addr := net.JoinHostPort(*httpHost, *httpPort)
		logger.Log("transport", "HTTP", "addr", addr)
		errs <- http.ListenAndServe(addr, handler)
	}()

	logger.Log("exit", <-errs)
}
