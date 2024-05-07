package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DmitriyKost/wbl0/pkg"
	_ "github.com/DmitriyKost/wbl0/pkg"
	"github.com/DmitriyKost/wbl0/pkg/database"
	"github.com/DmitriyKost/wbl0/pkg/structs"
	"github.com/nats-io/stan.go"
)


var sc stan.Conn

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sig)
		<-sig

		// Cancel the context to trigger graceful shutdown
		cancel()
	}()

	// Connecting to NATS streaming server
	clusterID := os.Getenv("CLUSTER_ID")
	clientID := os.Getenv("SUBSCRIBER_ID")
	ncURL := os.Getenv("NC_URL")
	channel := os.Getenv("CHANNEL")

	var err error
	sc, err = stan.Connect(clusterID, clientID, stan.NatsURL(ncURL))
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

	// Start the HTTP server
	router := http.NewServeMux()
	router.HandleFunc("/", pkg.Index)
	router.HandleFunc("/get/{order_uid}/", pkg.GetOrder)

    serverAddress := os.Getenv("SERVER_ADRESS")
	server := &http.Server{
		Addr:    serverAddress,
		Handler: pkg.LoggingMiddleware(router),
	}

	go func() {
		log.Printf("Listening on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for the context to be canceled
	<-ctx.Done()

	// Close the NATS subscription and connection
	subscription.Close()
	sc.Close()

	// Create a deadline for the shutdown process
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	// Shutdown the HTTP server gracefully
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	log.Println("Server gracefully stopped")
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
