package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) Serve() error {

	server := http.Server{
		Addr:         fmt.Sprintf(":%v", app.config.port),
		Handler:      app.Mount(),
		IdleTimeout:  time.Minute,
		WriteTimeout: 20 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	shutDownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
		s := <-quit

		log.Printf("Signal caught: %s, shutting down server", s)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		shutDownError <- server.Shutdown(ctx)
	}()

	log.Printf("Server starting at port %v\n", app.config.port)
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Wait for shutdown or timeout error
	if err := <-shutDownError; err != nil {
		return err
	}

	log.Println("Server stopped")
	return nil
}
