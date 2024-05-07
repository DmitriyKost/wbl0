package main

import (
	"log"
	"os"

	_ "github.com/DmitriyKost/wbl0/config"
	stan "github.com/nats-io/stan.go"
)

// Script to test NATS-streaming functionality, that publishes json from ./scripts/json/data.json file.
func main() {
    clusterID := os.Getenv("CLUSTER_ID")
    clientID := os.Getenv("PUBLISHER_ID")
    ncURL := os.Getenv("NC_URL")

    sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(ncURL))
    if err != nil {
        log.Fatalf("Failed to connect to NATS Streaming server: %v", err)
    }
    defer sc.Close()

    channel := os.Getenv("CHANNEL")

    data, err := os.ReadFile("./scripts/json/data.json")
    if err != nil {
        log.Fatalf("Failed to read JSON data: %v", err)
    }
    err = sc.Publish(channel, data)
    if err != nil {
        log.Fatalf("Failed to publish message: %v", err)
    } else {
        log.Println("Published JSON data succesfully")
    }
}
