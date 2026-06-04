package store

import (
	"fmt"
	"sync"
	"time"
)

type StoreInterface interface {
	Append(topic string, partition int, content string)
	Get(topic string, partition int, offset int) ([]Message, error)
	GetAllPerTopic(topic string) ([]Message, error)
	GetAll() ([]Message, error)

	CommitOffset(group string, topic string, partition int, offset int)
	GetOffset(group string, topic string, partition int) int 
}

type Message struct {
	Offset    int64
	Content   string
	CreatedAt time.Time
}

type Store struct {
	mu     sync.RWMutex
	topics map[string]map[int][]Message
	//M3
	offsets map[string]map[string]map[int]int
}

func NewStore() *Store {
	return &Store{
		topics: make(map[string]map[int][]Message),
		offsets: make(map[string]map[string]map[int]int),
	}
}

/********************************************************************/
func (s *Store) Append(topic string, partition int, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.topics[topic] == nil {
		s.topics[topic] = make(map[int][]Message)
	}

	s.topics[topic][partition] = append(
		s.topics[topic][partition],
		Message{
			Offset:    int64(len(s.topics[topic][partition])),
			Content:   content,
			CreatedAt: time.Now(),
		},
	)
}

/********************************************************************/
func (s *Store) Get(topic string, partition int, offset int) ([]Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.topics[topic] == nil || s.topics[topic][partition] == nil {
		return nil, fmt.Errorf("topic or partition not found")
	}

	answer := []Message{}

	for ; offset < len(s.topics[topic][partition]); offset++ {
		answer = append(
			answer,
			s.topics[topic][partition][offset],
		)
	}
	return answer, nil
}

/********************************************************************/
func (s *Store) GetAllPerTopic(topic string) ([]Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.topics[topic] == nil {
		return nil, fmt.Errorf("topic not found")
	}

	answer := []Message{}

	for _, messages := range s.topics[topic] {
		answer = append(answer, messages...)
	}

	return answer, nil
}

/********************************************************************/
func (s *Store) GetAll() ([]Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	answer := []Message{}

	for _, partitionsForTopic := range s.topics {
		for _, messages := range partitionsForTopic {
			answer = append(answer, messages...)
		}
	}

	return answer, nil
}
// M3
/********************************************************************/
func (s *Store) GetOffset(group string, topic string, partition int) int {

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.offsets[group] == nil || s.offsets[group][topic] == nil {
		return 0 //If the group or topic don't exist yet it's okay we just return 0
	}
	return s.offsets[group][topic][partition]
}

func (s *Store) CommitOffset(group string, topic string, partition int, offset int) {

	s.mu.Lock()
	defer s.mu.Unlock()
	
	//Just create it 
	if s.offsets[group] == nil {
		s.offsets[group] = make(map[string]map[int]int)
	}
	if s.offsets[group][topic] == nil {
		s.offsets[group][topic] = make(map[int]int)
	}

	s.offsets[group][topic][partition] = offset	
}

