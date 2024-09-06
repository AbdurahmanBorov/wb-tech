package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nats-io/stan.go"

	"wb-tech/internal/config"
	"wb-tech/internal/models"
)

func main() {
	file, err := os.ReadFile("model.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var order models.Order
	err = json.Unmarshal(file, &order)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	conn, err := stan.Connect(cfg.NatsConfig.ClusterID, "test-client2")
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer conn.Close()

	orderData, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Error marshalling order data: %v", err)
	}

	if err := conn.Publish(cfg.NatsConfig.NatsChan, orderData); err != nil {
		log.Fatalf("Error publishing message: %v", err)
	}

	log.Println("Message published successfully")
}
