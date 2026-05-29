package broker

import (
	"hash/fnv"

	"github.com/Xenios7/Raven/internal/store"
)

type BrokerInterface interface {
    Publish(content string, topic string, key string, numPartitions int)
    Consume(topic string, partition int, offset int) ([]store.Message, error)
	ConsumeAllPerTopic(topic string) ([]store.Message, error)
	ConsumeAll() ([]store.Message, error)
}

type Broker struct {
    store store.StoreInterface //Reference to store struct(interface)
}

// Constructor
// Needs store interface (it doesn't know or care if the caller passed in a real *Store, a fake test store, or a disk-based store. It just knows it can call Append and Get on it.)
// so it can initialize the store that it has
func NewBroker(s store.StoreInterface) *Broker {
	return &Broker{
		store: s,
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










































