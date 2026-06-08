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

    // K8s StatefulSet role detection via pod hostname
    if role == "" {
        hostname, _ := os.Hostname()
        if strings.HasSuffix(hostname, "-0") {
            role = "leader"
            replicaAddrs = "raven-1.raven:9001,raven-2.raven:9002"
        } else if strings.HasSuffix(hostname, "-1") {
            role = "replica"
            grpcPort = "9001"
        } else {
            role = "replica"
            grpcPort = "9002"
        }
    }

    if role == "replica" {
        s := store.NewStore()
        replServer := broker.NewReplicationServer(s)
        go startGRPCServer(replServer, ":"+grpcPort)
    } else {
        s := store.NewStore()
        addrs := strings.Split(replicaAddrs, ",")
        replicator := broker.NewReplicator(addrs)
        b := broker.NewBroker(s, replicator)
        h := api.NewHandler(b)
        r := api.NewRouter(h)

        go func() {
            fmt.Println("HTTP server listening on :8080")
            if err := http.ListenAndServe(":8080", r); err != nil {
                fmt.Println("HTTP server error:", err)
            }
        }()
    }

    select {}
}