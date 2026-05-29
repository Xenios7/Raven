package api

import "net/http"


func NewRouter(h *Handler) *http.ServeMux {
    mux := http.NewServeMux()
	
	mux.HandleFunc("POST /publish", h.PublishHandler)
	mux.HandleFunc("GET /consume", h.ConsumeHandler)

	return mux
}






















