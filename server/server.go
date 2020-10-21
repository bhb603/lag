package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

// Config holds server config
type Config struct {
	Port        string
	MaxLag      time.Duration
	MaxDataSize string
}

// Serve starts the server
func Serve(cfg *Config) {
	r := mux.NewRouter()

	// Middleware
	r.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		log.Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Interface("cfRay", r.Header[http.CanonicalHeaderKey("cf-ray")]).
			Msg("")
	}))
	r.Use(parseParamsMiddleware)
	lm := lagMiddleware{
		maxLag: cfg.MaxLag,
	}
	r.Use(lm.Middleware)

	// Routes
	r.HandleFunc("/", indexHandler)

	if dh, err := newDataHandler(cfg.MaxDataSize); err == nil {
		r.HandleFunc("/data", dh.Handler)
	} else {
		log.Fatal().Err(err).Send()
	}

	r.HandleFunc("/error/{status:[0-9]+}", errorHandler)
	r.HandleFunc("/headers", headersHandler)
	r.HandleFunc("/health", healthHandler)
	r.PathPrefix("/").Handler(http.NotFoundHandler())

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", cfg.Port),
		Handler:      r,
		WriteTimeout: cfg.MaxLag + 1*time.Second,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	go func() {
		log.Printf("listening on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil {
			log.Err(err).Send()
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Print("received interrupt signal")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)
	os.Exit(0)
}
