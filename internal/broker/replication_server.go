//server side. The replica runs this to receive messages from the leader.

package broker

import (
	"context"

	"github.com/Xenios7/Raven/internal/store"
	"github.com/Xenios7/Raven/proto"
)

type ReplicationServer struct {
	proto.UnimplementedReplicationServiceServer
	store store.StoreInterface
}

func NewReplicationServer(s store.StoreInterface) *ReplicationServer{
	return &ReplicationServer{
		store: s,
	}
}
// Replicate receives a message from the leader broker and stores it in the local store.
// Called automatically by the gRPC server when the leader replicates a message.
func (rs *ReplicationServer) Replicate(ctx context.Context, msg *proto.ReplicaMessage) (*proto.ReplicateResponse, error){

	//Store to local store struct
	rs.store.Append(msg.Topic, int(msg.Partition), msg.Content)

	//Send response back
	response := &proto.ReplicateResponse{
		Success: true,
	}
	return response, nil


}






