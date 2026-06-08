package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"log/slog"
	"pulse-parse/worker/internal/parser"
)

func main () {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	logger.Info("Starting PulseParse worker service", "port", ":8080")

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	mux.HandleFunc("/api/v1/parse", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req parser.LabReport
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			slog.Error("Failed to decode request body", "error", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		slog.Info("Processing lab report", "report_id", req.ID, "patient_id", req.PatientID)

		
		result, err := parser.ExtractAndAnalyze(req)
		if err != nil {
			slog.Error("Pipeline processing failure", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		logger.Info("Shutting down server dynamically", "signal", s.String())

		// Give outstanding requests 30 seconds to complete
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- server.Shutdown(ctx)
	}()

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Server failed to start or crash unexpectedly", "error", err)
		os.Exit(1)
	}

	err = <-shutdownError
	if err != nil {
		logger.Error("Error during graceful shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Worker service stopped cleanly")
}