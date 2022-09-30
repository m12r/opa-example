package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/m12r/opa-example/internal/server"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	addr := ":8080"
	if addrFromEnv := os.Getenv("OPA_EXAMPLE_ADDR"); addrFromEnv != "" {
		addr = addrFromEnv
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s := &http.Server{
		Handler: server.NewServer(),
	}
	//nolint:errcheck
	go s.Serve(l)

	<-ctx.Done()
	stopHTTPServer(s)
	return nil
}

func stopHTTPServer(s *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	s.SetKeepAlivesEnabled(false)
	if err := s.Shutdown(ctx); err != nil {
		log.Printf("warn: could not gracefully shutdown http server: %v", err)
	}
}
