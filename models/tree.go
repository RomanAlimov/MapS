package models

type Tree struct {
	ID           uint    `json:"id" gorm:"primaryKey"`
	Species      string  `json:"species" gorm:"not null"`
	Lat          float64 `json:"lat" gorm:"not null"`
	Lng          float64 `json:"lng" gorm:"not null"`
	PlantedYear  *int    `json:"planted_year"`
	HealthStatus string  `json:"health_status" gorm:"default:healthy"`
}
