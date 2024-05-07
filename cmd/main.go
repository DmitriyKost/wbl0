package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/DmitriyKost/wbl0/pkg"
	_ "github.com/DmitriyKost/wbl0/pkg"
	"github.com/DmitriyKost/wbl0/pkg/database"
	"github.com/DmitriyKost/wbl0/pkg/structs"
	"github.com/nats-io/stan.go"
)


func main() {
	clusterID := os.Getenv("CLUSTER_ID")
	clientID := os.Getenv("SUBSCRIBER_ID")
	ncURL := os.Getenv("NC_URL")
    channel := os.Getenv("CHANNEL")

	// Connecting to NATS streaming server
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(ncURL))
	if err != nil {
		log.Fatalf("Failed to connect to NATS Streaming server: %v", err)
	}
	defer sc.Close()

	// Subscribing to channel
	subscription, err := sc.Subscribe(channel, handleMsg, stan.DurableName("durable-name"))
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}
    defer subscription.Close()
	log.Println("Subscriber connected and waiting for messages...")

	router := http.NewServeMux()
    router.HandleFunc("/", pkg.Index)
    router.HandleFunc("/get/{order_uid}/", pkg.GetOrder)

    server := http.Server {
        Addr: ":8080",
        Handler: pkg.LoggingMiddleware(router),
    }

    log.Printf("Listening on %s \n", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Error on server %w", err)
	}

    select {}
}

// Handles incoming messages received from NATS Streaming.
//
// It validates the incoming message and inserts it into the database if it's valid.
func handleMsg(msg *stan.Msg) {
    if err := validateMsg(msg.Data); err != nil {
        log.Printf("Received an incorrect message: %v", err)
    } else {
        if err := database.InsertOrder(msg.Data); err != nil {
            log.Printf("Failed to insert message: %v", err)
        } else {
            log.Println("Message received!")
        }
    }
}

// Validates the format of an incoming message.
//
// It attempts to unmarshal the message into a structured order object.
func validateMsg(msg []byte) error {
    var order structs.Order
    if err := json.Unmarshal(msg, &order); err != nil {
        return err
    }
    return nil
}
