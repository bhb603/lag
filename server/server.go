package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	r.HandleFunc("/", indexHandler)

	if dh, err := newDataHandler(cfg.MaxDataSize); err == nil {
		r.HandleFunc("/data", dh.Handler)
	} else {
		log.Fatal(err)
	}

	r.HandleFunc("/error/{status:[0-9]+}", errorHandler)
	r.HandleFunc("/headers", headersHandler)
	r.HandleFunc("/health", healthHandler)

	r.Use(parseParamsMiddleware)
	lm := lagMiddleware{
		maxLag: cfg.MaxLag,
	}
	r.Use(lm.Middleware)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", cfg.Port),
		Handler:      handlers.CombinedLoggingHandler(os.Stdout, r),
		WriteTimeout: cfg.MaxLag + 1*time.Second,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	go func() {
		log.Printf("listening on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("received interrupt signal")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)
	os.Exit(0)
}
