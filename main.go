package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	addr                    string
	gracefulShutdownTimeout time.Duration
	maxTimeDelay            time.Duration
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "server listen address")
	flag.DurationVar(&gracefulShutdownTimeout, "graceful-timeout", time.Second*30, "duration the server will wait before cancelling idle connections in a graceful shutdown - e.g. 15s or 1m")
	flag.DurationVar(&maxTimeDelay, "max-time-delay", time.Second*30, "max allowable time delay")
	flag.Parse()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	statusCode, err := strconv.Atoi(vars["status"])
	if err != nil || statusCode < 400 || statusCode >= 600 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func parseParamsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		next.ServeHTTP(w, r)
	})
}

func timeDelayMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if t := r.Form.Get("t"); len(t) > 0 {
			d, err := time.ParseDuration(t)
			if err != nil {
				http.Error(w, "Invalid time parameter", http.StatusBadRequest)
				return
			}
			if d > maxTimeDelay {
				time.Sleep(maxTimeDelay)
			} else {
				time.Sleep(d)
			}
		}
		next.ServeHTTP(w, r)
	})

}

func main() {
	r := mux.NewRouter()
	loggedRouter := handlers.CombinedLoggingHandler(os.Stdout, r)

	// Routes
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/error/{status:[0-9]+}", errorHandler).Methods("GET")
	r.HandleFunc("/health", healthHandler)
	r.HandleFunc("/", notFoundHandler)
	r.Use(parseParamsMiddleware)
	r.Use(timeDelayMiddleware)

	srv := &http.Server{
		Addr:    addr,
		Handler: loggedRouter,
		// WriteTimeout: time.Second * 15,
		// ReadTimeout:  time.Second * 15,
		// IdleTimeout:  time.Second * 60,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		sig := <-sigint
		log.Printf("Recieved signal %s\n", sig)
		ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
		log.Printf("Idle connections closed")
	}()

	log.Printf("Server listening on %s", addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe: %v", err)
	} else {
		log.Printf("Server closed")
	}

	// block until idleConnsClosed
	<-idleConnsClosed
}
