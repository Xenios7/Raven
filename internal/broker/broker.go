package broker

import (
	"fmt"
	"hash/fnv"
	"time"

	"github.com/Xenios7/Raven/internal/store"
	"github.com/Xenios7/Raven/proto"
)

type BrokerInterface interface {
    Publish(content string, topic string, key string, numPartitions int)
    Consume(topic string, partition int, offset int) ([]store.Message, error)
	ConsumeAllPerTopic(topic string) ([]store.Message, error)
	ConsumeAll() ([]store.Message, error)

	ConsumeWithGroup(group string, topic string, partition int) ([]store.Message, error)

	Ack(group string, topic string, partition int, offset int)
}

type Broker struct {
    store store.StoreInterface //Reference to store struct(interface)
	//Replicator needed for updating replicas
	replicator ReplicatorInterface
}

// Constructor
// Needs store interface (it doesn't know or care if the caller passed in a real *Store, a fake test store, or a disk-based store. It just knows it can call Append and Get on it.)
// so it can initialize the store that it has
func NewBroker(s store.StoreInterface, r ReplicatorInterface) *Broker {
	return &Broker{
		store: s,
		replicator: r,
	}
}

func (b *Broker) partition(key string, numPartitions int) int{
	h := fnv.New32a() //create a hashing machine (convert any string into a number)
	h.Write([]byte(key)) //feed the key into it (strings must be converted to bytes first)
	return int(h.Sum32()) % numPartitions //get the result and fit it into the number of partitions using modulo

	/*Hash produces some big number, say 2847291
	2847291 % 3 = 0
	Message goes to partition 0
	*/
}

/********************************************************************/
func (b *Broker) Publish(content string, topic string, key string, numPartitions int){

	//Find partition
	partitionNumber := b.partition(key, numPartitions)

	//Append using store structure from store/
	b.store.Append(topic, partitionNumber, content)

	/*******************************/
	/*******************************/
	// We need to update the replicas
	msg := &proto.ReplicaMessage{ //No need for a constructor here, in Go we can do this in this way if we want, ofcourse constructor works as well but here we don't have one.
		Topic:     topic,
		Content:   content,
		Partition: int32(partitionNumber),
		//Offset will be calculated by the other broker node inside of store when it Appneds it. So we don't set it here.
		CreatedAt: time.Now().String(),
	}

	// b.replicator.Replicate(msg)
	if err := b.replicator.Replicate(msg); err != nil {
		fmt.Println("replication error:", err)
	}
}

/********************************************************************/
func (b *Broker) Consume(topic string, partition int, offset int) ([]store.Message, error){
	return b.store.Get(topic, partition, offset)
}

func (b *Broker) ConsumeAllPerTopic(topic string) ([]store.Message, error){
	return b.store.GetAllPerTopic(topic)
}

func (b *Broker) ConsumeAll() ([]store.Message, error){
	return b.store.GetAll()
}


/******************************* M3 *************************************/
//Does the same thing as Consume only difference is that it tracks offset per group, also we update offset after we proccess messages
func (b *Broker) ConsumeWithGroup(group string, topic string, partition int) ([]store.Message, error) {

	offset := b.store.GetOffset(group, topic, partition)

	messages, err := b.store.Get(topic, partition, offset)
	if err != nil {
		return nil, err
	}

	// newOffset := offset + len(messages) // New offset = old offset + number of messages prossed

	// b.store.CommitOffset(group, topic, partition, newOffset)

	return  messages, err
}

func (b *Broker) Ack(group string, topic string, partition int, offset int) {
    b.store.CommitOffset(group, topic, partition, offset)
}







































