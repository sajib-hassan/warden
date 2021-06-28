package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

// Server provides an http.Server.
type Server struct {
	*http.Server
}

// NewServer creates and configures an APIServer serving all application routes.
func NewServer() (*Server, error) {
	log.Println("configuring server...")
	api, err := New()
	if err != nil {
		return nil, err
	}

	port := viper.GetString("port")
	srv := http.Server{
		ReadTimeout:  viper.GetDuration("READ_TIMEOUT") * time.Second,
		WriteTimeout: viper.GetDuration("WRITE_TIMEOUT") * time.Second,
		IdleTimeout:  viper.GetDuration("IDLE_TIMEOUT") * time.Second,
		Addr:         ":" + port,
		Handler:      api,
	}

	return &Server{&srv}, nil
}

// Start runs ListenAndServe on the http.Server with graceful shutdown.
func (srv *Server) Start() {
	log.Println("starting server...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		log.Println("Server Listening on :" + srv.Addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	sig := <-quit
	log.Println("Shutting down server... Reason:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Println("Server gracefully stopped")
}
