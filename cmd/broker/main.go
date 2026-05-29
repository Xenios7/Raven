package main

import (
	"fmt"
	"net/http"

	"github.com/Xenios7/Raven/internal/api"
	"github.com/Xenios7/Raven/internal/broker"
	"github.com/Xenios7/Raven/internal/store"
)

func main() {
    fmt.Println("Raven broker starting...")
    
    // wire together store → broker → handler → router
    s := store.NewStore()
    b := broker.NewBroker(s)
    h := api.NewHandler(b)
    r := api.NewRouter(h)

    // start listening for incoming HTTP requests
    if err := http.ListenAndServe(":8080", r); err != nil {
        fmt.Println("server error:", err)
    }


}














