package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"website/backend/internal/api"
	"website/backend/internal/database"
	"website/backend/internal/mqtt"
	"website/backend/internal/ws"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Initialize Database
	database.InitDB()

	// Initialize MQTT
	mqtt.InitMQTT()

	// Setup Router
	r := mux.NewRouter()

	// Public API Routes
	r.HandleFunc("/api/nodes", api.GetNodes).Methods("GET")
	r.HandleFunc("/api/nodes/{mac}/history", api.GetSensorHistory).Methods("GET")
	r.HandleFunc("/api/harvests", api.GetHarvests).Methods("GET")
	r.HandleFunc("/api/harvests/averages", api.GetHarvestSoilAverages).Methods("GET")
	r.HandleFunc("/api/news", api.GetEconomicNews).Methods("GET")
	r.HandleFunc("/api/cpo-price", api.GetCPOPrice).Methods("GET")
	r.HandleFunc("/api/auth/login", api.Login).Methods("POST")
	r.HandleFunc("/api/auth/status", api.CheckAuth).Methods("GET", "OPTIONS")

	// Protected Admin Routes
	adminRouter := r.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(api.AuthMiddleware)
	adminRouter.HandleFunc("/nodes/{mac}", api.UpdateNodeAlias).Methods("PUT", "OPTIONS")
	adminRouter.HandleFunc("/harvests", api.CreateHarvest).Methods("POST", "OPTIONS")

	// WebSocket Route
	r.HandleFunc("/ws", ws.HandleWebSocket)

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Adjust for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
