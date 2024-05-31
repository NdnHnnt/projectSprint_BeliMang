package models

import (
	"time"
)

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

type MerchantItems struct {
	ID              string    `json:"merchant_id" db:"id"`
	Name            string    `json:"name" db:"name"`
	ProductCategory string    `json:"productCategory" db:"productCategory"`
	ImageUrl        string    `json:"imageUrl" db:"imageUrl"`
	Price           int       `json:"price" db:"price"`
	MerchantId      string    `json:"merchantId" db:"merchantId"`
	CreatedAt       time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updatedAt"`
}

type MerchantEstimatePrice struct {
	UserLocation struct {
		Lat  float64 `json:"lat"`
		Long float64 `json:"long"`
	} `json:"userLocation"`
	Orders []MerchantEstimateOrder `json:"orders"`
}

type MerchantEstimateOrder struct {
	MerchantId      string                      `json:"merchantId"`
	IsStartingPoint bool                        `json:"isStartingPoint"`
	Items           []MerchantEstimateOrderItem `json:"items"`
}

type MerchantEstimateOrderItem struct {
	ItemId   string `json:"itemId"`
	Quantity int    `json:"quantity"`
}
