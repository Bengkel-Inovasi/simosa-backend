package models

import (
	"time"
)

type Node struct {
	MacAddress   string    `gorm:"primaryKey" json:"mac_address"`
	Alias        string    `json:"alias"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	RegisteredAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"registered_at"`
	IsRegistered bool      `gorm:"default:false" json:"is_registered"`
}

type SensorReading struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	NodeMac     string    `json:"node_mac"`
	PH          float32   `json:"ph"`
	N           float32   `json:"n"`
	P           float32   `json:"p"`
	K           float32   `json:"k"`
	Moisture    float32   `json:"moisture"`
	Temperature float32   `json:"temperature"`
	Timestamp   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"timestamp"`
}

type Harvest struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	HarvestDate time.Time `json:"harvest_date"`
	YieldKg     float32   `json:"yield_kg"`
	PricePerKg  float32   `json:"price_per_kg"`
	GrossIncome float32   `json:"gross_income"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	Expenses    []Expense `gorm:"foreignKey:HarvestID" json:"expenses"`
}

type Expense struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	HarvestID uint    `json:"harvest_id"`
	Name      string  `json:"name"`
	Amount    float32 `json:"amount"`
}
