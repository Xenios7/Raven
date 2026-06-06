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

type AckRequest struct {	
	Group 		 	string `json:"group"`
	Topic			string `json:"topic"`
	Partition		int    `json:"partition"`
	Offset			int    `json:"offset"`
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

func (h *Handler) ConsumeAllPerTopicHandler(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		http.Error(w, "topic is required", http.StatusBadRequest)
		return
	}

	messages, err := h.broker.ConsumeAllPerTopic(topic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *Handler) ConsumeAllHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := h.broker.ConsumeAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *Handler) ConsumeWithGroupHandler(w http.ResponseWriter, r *http.Request) {

	topic := r.URL.Query().Get("topic")
	if topic == "" {
		http.Error(w, "topic is required", http.StatusBadRequest)
   		return
	}
	
	partition, err := strconv.Atoi(r.URL.Query().Get("partition"))
	if err != nil {
		http.Error(w, "invalid partition", http.StatusBadRequest)
   		return
	}

	//Get the group name (not offset like before)
	group := r.URL.Query().Get("group")
	if group == "" {
		http.Error(w, "group is required", http.StatusBadRequest)
   		return
	}

	messages, err := h.broker.ConsumeWithGroup(group, topic, partition)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
    	return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// The ack endpoint is the explicit "I finished processing" signal.
// Without it the offset never moves and the consumer keeps getting the same messages.
func (h *Handler) AckHandler(w http.ResponseWriter, r *http.Request) {

	var req AckRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
        return
	}

	h.broker.Ack(req.Group, req.Topic, req.Partition, req.Offset)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("offset committed"))

}

























