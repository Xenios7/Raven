package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

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

    role := os.Getenv("ROLE")
    grpcPort := os.Getenv("GRPC_PORT")
    replicaAddrs := os.Getenv("REPLICA_ADDRS")

    if role == "replica" {
        s := store.NewStore()
        replServer := broker.NewReplicationServer(s)
		// runs forever waiting for gRPC calls
        go startGRPCServer(replServer, ":"+grpcPort) //use goroutine because startGRPCServer blocks forever, so it runs in the background.
    } else {
        // leader
        s := store.NewStore()
        addrs := strings.Split(replicaAddrs, ",")
        replicator := broker.NewReplicator(addrs)
        b := broker.NewBroker(s, replicator)
        h := api.NewHandler(b)
        r := api.NewRouter(h)

		// runs forever waiting for http requests
        go func() { // use goroutine http.ListenAndServe never stops it runs forever waiting for requests
            fmt.Println("HTTP server listening on :8080")
            if err := http.ListenAndServe(":8080", r); err != nil {
                fmt.Println("HTTP server error:", err)
            }
        }()
    }

    select {}
}