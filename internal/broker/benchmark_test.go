package broker

import (
	"github.com/Xenios7/Raven/internal/store"
	"testing"
)

func BenchmarkPublish(b *testing.B) {
	s := store.NewStore()
	r := NewFakeReplicator([]store.StoreInterface{})
	br := NewBroker(s, r)

	for i := 0; i < b.N; i++ {
		br.Publish("video1", "uploads", "user-42", 3)
	}
}

func BenchmarkConsume(b *testing.B) {
	s := store.NewStore()
	r := NewFakeReplicator([]store.StoreInterface{})
	br := NewBroker(s, r)

	br.Publish("video1", "uploads", "user-42", 3)
	partitionNum := br.partition("user-42", 3)

	for i := 0; i < b.N; i++ {
		br.ConsumeWithGroup("billing", "uploads", partitionNum)
	}
}