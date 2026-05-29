package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Xenios7/Raven/internal/broker"
)


type Handler struct {
	broker broker.BrokerInterface
}

func NewHandler(b broker.BrokerInterface) *Handler{
	return &Handler{
		broker: b,
	}
}

type PublishRequest struct {
    Content       string `json:"content"`
    Topic         string `json:"topic"`
    Key           string `json:"key"`
    NumPartitions int    `json:"numPartitions"`
}

func (h *Handler) PublishHandler(w http.ResponseWriter, r *http.Request) {

	var req PublishRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
        return
	}

	h.broker.Publish(req.Content, req.Topic, req.Key, req.NumPartitions)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("message published"))
}

func (h *Handler) ConsumeHandler(w http.ResponseWriter, r *http.Request) {

	// topic := r.URL.Query().Get("topic")
	// offset := r.URL.Query().Get("offset")
	// partition := r.URL.Query().Get("partition")

	topic := r.URL.Query().Get("topic")
	if topic == "" {
		http.Error(w, "topic is required", http.StatusBadRequest)
   		return
	}
	
	//Convert string to integer also check that offset is valid
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		http.Error(w, "invalid offset", http.StatusBadRequest)
   		return
	}

	//Convert string to integer also check that partition is valid
	partition, err := strconv.Atoi(r.URL.Query().Get("partition"))
	if err != nil {
		http.Error(w, "invalid partition", http.StatusBadRequest)
   		return
	}

	messages, err := h.broker.Consume(topic, partition, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
    	return
	}

	//Convert messages (structs) to JSON and write them back
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}





























