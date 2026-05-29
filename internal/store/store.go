package store


import (
    "fmt"
    "sync"
    "time"
)

type StoreInterface interface {
    Append(topic string, partition int, content string)
    Get(topic string, partition int, offset int) ([]Message, error)
}

type Message struct {
    Offset    int64
    Content   string
    CreatedAt time.Time
}

type Store struct {
    mu         sync.RWMutex
    partitions map[string]map[int][]Message
}

// Constructor
func NewStore() *Store {
	return &Store{
        partitions: make(map[string]map[int][]Message),		
	}
}

/********************************************************************/
func (s *Store) Append(topic string, partition int, content string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if s.partitions[topic] == nil {
        s.partitions[topic] = make(map[int][]Message)
    }

    s.partitions[topic][partition] = append(
        s.partitions[topic][partition],
        Message{
            Offset:    int64(len(s.partitions[topic][partition])),
            Content:   content,
            CreatedAt: time.Now(),
        },
    )
}
/********************************************************************/
func (s *Store) Get(topic string, partition int, offset int) ([]Message, error) {

	s.mu.RLock() 
	defer s.mu.RUnlock() 

	if s.partitions[topic] == nil || s.partitions[topic][partition] == nil  {
		return nil, fmt.Errorf("topic or partition not found")
	}
	
	answer := []Message{} //same as var answer []Message

	for ; offset < len(s.partitions[topic][partition]); offset++ {
		answer = append(
			answer, 
			s.partitions[topic][partition][offset],
		)
	}
	return answer, nil
}










