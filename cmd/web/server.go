package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/spankie/tw-interview/blockparser"
)

type Server struct {
	parser blockparser.BlockParser
}

type response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   string `json:"error"`
}

func respond(w http.ResponseWriter, status int, res any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error(fmt.Sprintf("error encoding response: %v", err))
	}
}

func newServer(blockParser blockparser.BlockParser) *http.Server {
	server := &Server{
		parser: blockParser,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /block", server.getCurrentBlockNumber)
	mux.HandleFunc("GET /transactions/{address}", server.getTransactionsByAddress)
	mux.HandleFunc("GET /subscribe/{address}", server.subscribeToAddress)

	port := os.Getenv("TW_PORT")
	if port == "" {
		port = "8080"
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func (s *Server) getCurrentBlockNumber(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusOK, response{
		Message: "success",
		Data:    s.parser.GetCurrentBlock(),
		Error:   "",
	})
}

func (s *Server) getTransactionsByAddress(w http.ResponseWriter, r *http.Request) {
	address := r.PathValue("address")
	respond(w, http.StatusOK, response{
		Message: "success",
		Data:    s.parser.GetTransactions(address),
		Error:   "",
	})
}

func (s *Server) subscribeToAddress(w http.ResponseWriter, r *http.Request) {
	address := r.PathValue("address")
	if !s.parser.Subscribe(address) {
		respond(w, http.StatusBadRequest, response{
			Message: "",
			Data:    "",
			Error:   "address could not be subscribed",
		})
		return
	}

	respond(w, http.StatusOK, response{
		Message: fmt.Sprintf("Subscribed to address: %s", address),
		Data:    "",
		Error:   "",
	})
}
