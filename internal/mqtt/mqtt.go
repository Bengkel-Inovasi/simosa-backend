package mqtt

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"website/backend/internal/database"
	"website/backend/internal/models"
	"website/backend/internal/ws"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var Client mqtt.Client

type SensorPayload struct {
	MacAddress  string  `json:"mac_address"`
	PH          float32 `json:"ph"`
	N           float32 `json:"n"`
	P           float32 `json:"p"`
	K           float32 `json:"k"`
	Moisture    float32 `json:"moisture"`
	Temperature float32 `json:"temperature"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

func InitMQTT() {
	opts := mqtt.NewClientOptions()
	
	host := os.Getenv("MQTT_HOST")
	port := os.Getenv("MQTT_PORT")
	if port == "" {
		port = "1883"
	}

	scheme := "tcp"
	if port == "8883" || port == "8884" || os.Getenv("MQTT_USE_TLS") == "true" {
		scheme = "ssl"
		opts.SetTLSConfig(&tls.Config{
			InsecureSkipVerify: false,
		})
	}

	brokerURL := fmt.Sprintf("%s://%s:%s", scheme, host, port)
	fmt.Printf("Connecting to MQTT Broker: %s...\n", brokerURL)
	opts.AddBroker(brokerURL)

	opts.SetUsername(os.Getenv("MQTT_USER"))
	opts.SetPassword(os.Getenv("MQTT_PASSWORD"))
	opts.SetClientID("simosa_backend")

	opts.OnConnect = func(c mqtt.Client) {
		fmt.Printf("Connected to MQTT Broker: %s\n", brokerURL)
		if token := c.Subscribe("simosa/nodes/+/data", 0, onMessageReceived); token.Wait() && token.Error() != nil {
			log.Fatalf("Failed to subscribe: %v", token.Error())
		}
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Printf("Connection lost: %v\n", err)
	}

	Client = mqtt.NewClient(opts)
	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT: %v", token.Error())
	}
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	var payload SensorPayload
	err := json.Unmarshal(message.Payload(), &payload)
	if err != nil {
		fmt.Printf("Error unmarshaling MQTT payload: %v\n", err)
		return
	}

	// 1. Check if node exists
	var node models.Node
	result := database.DB.First(&node, "mac_address = ?", payload.MacAddress)
	if result.Error != nil {
		// New node detected
		node = models.Node{
			MacAddress:   payload.MacAddress,
			Latitude:     payload.Latitude,
			Longitude:    payload.Longitude,
			IsRegistered: false,
		}
		database.DB.Create(&node)
		
		// Notify frontend about new unregistered node
		ws.BroadcastNotification("new unregistered nodes detected")
	}

	// 2. Save reading
	reading := models.SensorReading{
		NodeMac:     payload.MacAddress,
		PH:          payload.PH,
		N:           payload.N,
		P:           payload.P,
		K:           payload.K,
		Moisture:    payload.Moisture,
		Temperature: payload.Temperature,
	}
	database.DB.Create(&reading)

	// 3. Broadcast to frontend via WebSocket
	ws.BroadcastReading(reading)
}
