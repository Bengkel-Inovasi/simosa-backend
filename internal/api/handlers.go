package api

import (
	"encoding/json"
	"net/http"
	"website/backend/internal/database"
	"website/backend/internal/models"

	"github.com/gorilla/mux"
)

func GetNodes(w http.ResponseWriter, r *http.Request) {
	var nodes []models.Node
	database.DB.Find(&nodes)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

func UpdateNodeAlias(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	mac := params["mac"]

	var node models.Node
	if err := database.DB.First(&node, "mac_address = ?", mac).Error; err != nil {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	var update struct {
		Alias        string `json:"alias"`
		IsRegistered bool   `json:"is_registered"`
	}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	node.Alias = update.Alias
	node.IsRegistered = update.IsRegistered
	database.DB.Save(&node)

	json.NewEncoder(w).Encode(node)
}

func GetSensorHistory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	mac := params["mac"]

	var readings []models.SensorReading
	database.DB.Where("node_mac = ?", mac).Order("timestamp desc").Limit(100).Find(&readings)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(readings)
}

func CreateHarvest(w http.ResponseWriter, r *http.Request) {
	var harvest models.Harvest
	if err := json.NewDecoder(r.Body).Decode(&harvest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	harvest.GrossIncome = harvest.YieldKg * harvest.PricePerKg
	database.DB.Create(&harvest)

	json.NewEncoder(w).Encode(harvest)
}

func GetHarvests(w http.ResponseWriter, r *http.Request) {
	var harvests []models.Harvest
	database.DB.Preload("Expenses").Find(&harvests)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(harvests)
}

func GetHarvestSoilAverages(w http.ResponseWriter, r *http.Request) {
	// Implementation for calculating average soil data between harvests
	// This would take two dates and average the sensor_readings in that range
	// For brevity, I'll implement a basic version
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	var result struct {
		AvgPH          float64 `json:"avg_ph"`
		AvgN           float64 `json:"avg_n"`
		AvgP           float64 `json:"avg_p"`
		AvgK           float64 `json:"avg_k"`
		AvgMoisture    float64 `json:"avg_moisture"`
		AvgTemperature float64 `json:"avg_temperature"`
	}

	database.DB.Model(&models.SensorReading{}).
		Select("AVG(ph) as avg_ph, AVG(n) as avg_n, AVG(p) as avg_p, AVG(k) as avg_k, AVG(moisture) as avg_moisture, AVG(temperature) as avg_temperature").
		Where("timestamp BETWEEN ? AND ?", start, end).
		Scan(&result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
