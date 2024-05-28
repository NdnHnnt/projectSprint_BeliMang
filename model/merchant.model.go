package models

import (
	"time"
)

type MerchantModel struct {
	ID               string    `json:"merchant_id" db:"id"`
	Name             string    `json:"name" db:"name"`
	MerchantCategory string    `json:"merchantCategory" db:"merchantCategory"`
	ImageUrl         string    `json:"imageUrl" db:"imageUrl"`
	Lat              float64   `json:"lat" db:"lat"`
	Lon              float64   `json:"lon" db:"lon"`
	CreatedAt        time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt" db:"updatedAt"`
}

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type MerchantRequest struct {
	Name             string   `json:"name"`
	MerchantCategory string   `json:"merchantCategory"`
	ImageUrl         string   `json:"imageUrl"`
	Location         Location `json:"location"`
}
