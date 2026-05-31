package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/Xenios7/Raven/internal/api"
	"github.com/Xenios7/Raven/internal/broker"
	"github.com/Xenios7/Raven/internal/store"
	"github.com/Xenios7/Raven/proto"
	"google.golang.org/grpc"
)

// startGRPCServer starts a gRPC replication server on the given address
func startGRPCServer(replServer *broker.ReplicationServer, addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("gRPC listen error:", err)
		return
	}
	grpcServer := grpc.NewServer()
	proto.RegisterReplicationServiceServer(grpcServer, replServer)
	fmt.Println("gRPC server listening on", addr)
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Println("gRPC server error:", err)
	}
}

func main() {
	fmt.Println("Raven broker starting...")

	// --- Replica 2 (gRPC only, no HTTP) ---
	s2 := store.NewStore()
	replServer2 := broker.NewReplicationServer(s2)
	go startGRPCServer(replServer2, ":9001")

	// --- Replica 3 (gRPC only, no HTTP) ---
	s3 := store.NewStore()
	replServer3 := broker.NewReplicationServer(s3)
	go startGRPCServer(replServer3, ":9002")

	// --- Leader Broker (HTTP + replicates to :9001 and :9002) ---
	s1 := store.NewStore()
	replicator := broker.NewReplicator([]string{":9001", ":9002"})
	b := broker.NewBroker(s1, replicator)
	h := api.NewHandler(b)
	r := api.NewRouter(h)

	// start HTTP server in background goroutine
	go func() {
		fmt.Println("HTTP server listening on :8080")
		if err := http.ListenAndServe(":8080", r); err != nil {
			fmt.Println("HTTP server error:", err)
		}
	}()

	// block main from exiting — keep all goroutines alive
	select {}
}