//client side. The leader calls this to send messages to replicas.
package broker

import (
    "context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    proto "github.com/Xenios7/Raven/proto"
)

type ReplicatorInterface interface {
    Replicate(msg *proto.ReplicaMessage) error
}

type Replicator struct {
	replicaAddresses []string
}

func NewReplicator(addresses []string) *Replicator {
    return &Replicator{
        replicaAddresses: addresses,
    }
}


func (r *Replicator) Replicate(msg *proto.ReplicaMessage) error {

	for _, address := range r.replicaAddresses {

		// 1. Create a connection to the replica
		conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		defer conn.Close()

		// 2. Create a client using the generated code
		client := proto.NewReplicationServiceClient(conn)

		// 3. Call Replicate
		_, err = client.Replicate(context.Background(), msg)
		if err != nil {
			return err
		}

	}

	return nil
}