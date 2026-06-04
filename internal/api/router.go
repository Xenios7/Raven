package api

import "net/http"

func NewRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /publish", h.PublishHandler)
	mux.HandleFunc("GET /consume", h.ConsumeHandler)
	mux.HandleFunc("GET /consume/topic", h.ConsumeAllPerTopicHandler)
	mux.HandleFunc("GET /consume/all", h.ConsumeAllHandler)
	mux.HandleFunc("GET /consume/group", h.ConsumeWithGroupHandler)
	
	return mux
}