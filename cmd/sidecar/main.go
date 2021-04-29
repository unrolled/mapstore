package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/unrolled/slog"
)

func main() {
	// Generate the router.
	mux := http.NewServeMux()

	healthEndpoint := "/.internal/healthz"
	if ep, ok := os.LookupEnv("HEALTH_ENDPOINT"); ok {
		healthEndpoint = ep
	}

	mux.Handle(healthEndpoint, healthz())

	// Check if we should enabled basic auth.
	baUser := strings.TrimSpace(os.Getenv("BASIC_AUTH_USER"))
	baPassword := strings.TrimSpace(os.Getenv("BASIC_AUTH_PASSWORD"))

	if len(baUser) != 0 && len(baPassword) != 0 {
		mux.Handle("/", basicAuth(api(), baUser, baPassword))
	} else {
		mux.Handle("/", api())
	}

	// Create the serving address.
	host := os.Getenv("HOST")
	port := "8080"
	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}
	addr := host + ":" + port

	// Start the server.
	server := &http.Server{
		Addr:    addr,
		Handler: log(recovery(noCache(slog.Requestify(mux)))),
	}

	// Start the server listening for connections.
	slog.Info("starting http server", slog.String("addr", addr))
	go func() {
		slog.Info("http server result", slog.Err(server.ListenAndServe()))
	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// Waiting for signal.
	<-stop

	// Create a timeout context that will shutdown the http server if needed.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start the server and log the final result.
	slog.Info("http server shutdown", slog.Err(server.Shutdown(ctx)))
}
