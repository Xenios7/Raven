package api

import (
	"net/http"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /publish", h.PublishHandler)
	mux.HandleFunc("GET /consume", h.ConsumeHandler)
	mux.HandleFunc("GET /consume/topic", h.ConsumeAllPerTopicHandler)
	mux.HandleFunc("GET /consume/all", h.ConsumeAllHandler)
	mux.HandleFunc("GET /consume/group", h.ConsumeWithGroupHandler)
	mux.HandleFunc("POST /ack", h.AckHandler)
    mux.Handle("/metrics", promhttp.Handler())

	return mux
}