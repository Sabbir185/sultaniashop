package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sabbir185/sultaniashop/internal/config"
	"github.com/Sabbir185/sultaniashop/internal/db"
	"github.com/Sabbir185/sultaniashop/internal/handlers"
	"github.com/Sabbir185/sultaniashop/internal/middleware"
)

func main() {
	cnf := config.LoadConfig()

	db, err := db.Connect(cnf.DB.Url)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to database successfully")

	loggerHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	logger := slog.New(loggerHandler)
	slog.SetDefault(logger)

	listingHandler := handlers.NewListingHandler(db, logger)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /list", listingHandler.ListAllProducts)
	mux.HandleFunc("DELETE /delete/{id}", listingHandler.DeleteList)
	mux.HandleFunc("POST /listings", listingHandler.DeleteList)

	handler := middleware.RequestId(mux)

	server := &http.Server{
		Addr:         ":" + cnf.App.Port,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	// Graceful Shutdown, error handling or listen and serve
	go func() {
		log.Printf("Starting server on port %s", cnf.App.Port)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to start: %v", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Printf("Server is shutting down with signal: %v\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown:\n%v\n", err.Error())
	}

	log.Println("Server exited properly")
}
