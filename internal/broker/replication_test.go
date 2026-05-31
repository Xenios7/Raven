package broker

import (
	"testing"

	"github.com/Xenios7/Raven/internal/store"
	proto "github.com/Xenios7/Raven/proto"
)

type FakeReplicator struct {
	replicaStores []store.StoreInterface
}

func NewFakeReplicator(r []store.StoreInterface) *FakeReplicator {
    return &FakeReplicator{
        replicaStores: r,
    }
}

// Loops through replica stores and appends the message directly.
// Bypasses gRPC — used for unit testing only.
func (f *FakeReplicator) Replicate(msg *proto.ReplicaMessage) error {

	for _, replica := range f.replicaStores {
		replica.Append(msg.Topic, int(msg.Partition), msg.Content)
	}
	return nil
}

func TestReplication(t *testing.T) {

	s1 := store.NewStore()
	s2 := store.NewStore()
	s3 := store.NewStore()

	freplicator := NewFakeReplicator([]store.StoreInterface{s2, s3})
	b := NewBroker(s1, freplicator)
	partitionNum := b.partition("user-1", 3)

	b.Publish("video123", "upload", "user-1", 3)

	// check s1 (leader)
	msgs1, err := s1.Get("upload", partitionNum, 0)
	if err != nil {
		t.Fatalf("s1 error: %v", err)
	}
	if len(msgs1) == 0 {
		t.Fatalf("s1 has no messages")
	}
	if msgs1[0].Content != "video123" {
		t.Fatalf("s1 wrong content: %s", msgs1[0].Content)
	}

	// check s2 (replica 1)
	msgs2, err := s2.Get("upload", partitionNum, 0)
	if err != nil {
		t.Fatalf("s2 error: %v", err)
	}
	if len(msgs2) == 0 {
		t.Fatalf("s2 has no messages")
	}
	if msgs2[0].Content != "video123" {
		t.Fatalf("s2 wrong content: %s", msgs2[0].Content)
	}

	// check s3 (replica 2)
	msgs3, err := s3.Get("upload", partitionNum, 0)
	if err != nil {
		t.Fatalf("s3 error: %v", err)
	}
	if len(msgs3) == 0 {
		t.Fatalf("s3 has no messages")
	}
	if msgs3[0].Content != "video123" {
		t.Fatalf("s3 wrong content: %s", msgs3[0].Content)
	}
}
