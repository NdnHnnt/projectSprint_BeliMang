package models

import (
	"time"
)

type MerchantCategory string

const (
	SmallRestaurant       MerchantCategory = "SmallRestaurant"
	MediumRestaurant      MerchantCategory = "MediumRestaurant"
	LargeRestaurant       MerchantCategory = "LargeRestaurant"
	MerchandiseRestaurant MerchantCategory = "MerchandiseRestaurant"
	BoothKiosk            MerchantCategory = "BoothKiosk"
	ConvenienceStore      MerchantCategory = "ConvenienceStore"
)

type ProductCategory string

const (
	Beverage   ProductCategory = "Beverage"
	Food       ProductCategory = "Food"
	Snack      ProductCategory = "Snack"
	Condiments ProductCategory = "Condiments"
	Additions  ProductCategory = "Additions"
)

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type MerchantRequest struct {
	Name             string           `json:"name"`
	MerchantCategory MerchantCategory `json:"merchantCategory"`
	ImageUrl         string           `json:"imageUrl"`
	Location         Location         `json:"location"`
}

type MerchantModel struct {
	ID               string           `json:"merchant_id" db:"id"`
	Name             string           `json:"name" db:"name"`
	MerchantCategory MerchantCategory `json:"merchantCategory" db:"merchantCategory"`
	ImageUrl         string           `json:"imageUrl" db:"imageUrl"`
	Lat              float64          `json:"lat" db:"lat"`
	Lon              float64          `json:"lon" db:"lon"`
	CreatedAt        time.Time        `json:"createdAt" db:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt" db:"updatedAt"`
}

type MerchantItems struct {
	ID              string          `json:"merchant_id" db:"id"`
	Name            string          `json:"name" db:"name"`
	ProductCategory ProductCategory `json:"productCategory" db:"productCategory"`
	ImageUrl        string          `json:"imageUrl" db:"imageUrl"`
	Price           int             `json:"price" db:"price"`
	MerchantId      string          `json:"merchantId" db:"merchantId"`
	CreatedAt       time.Time       `json:"createdAt" db:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt" db:"updatedAt"`
}
