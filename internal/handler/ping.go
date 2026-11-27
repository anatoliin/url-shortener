package handler

import (
	"context"
	"net/http"
)

type pinger interface {
	Ping(context.Context) error
}

type PingHandler struct {
	connection pinger
}

func NewPingHandler(connection pinger) *PingHandler {
	return &PingHandler{connection: connection}
}

func (h *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.connection.Ping(context.Background()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
