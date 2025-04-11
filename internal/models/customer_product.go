package models

import "time"

type Customer struct {
	CustomerID string `gorm:"primaryKey;type:varchar(50)"`
	Name       string
	Email      string
	Address    string
}

type Product struct {
	ProductID string `gorm:"primaryKey;type:varchar(50)"`
	Name      string
	Category  string
}

type Order struct {
	OrderID      string `gorm:"primaryKey;type:varchar(50)"`
	CustomerID   string `gorm:"type:varchar(50);not null"`
	Region       string
	DateOfSale   time.Time
	PaymentType  string
	ShippingCost float64
}

type OrderItem struct {
	ID           uint   `gorm:"primaryKey"`
	OrderID      string `gorm:"type:varchar(50);not null"`
	ProductID    string `gorm:"type:varchar(50);not null"`
	QuantitySold int
	UnitPrice    float64
	Discount     float64
}

type RevenueResult struct {
	ProductName string  `json:"product_name"`
	Revenue     float64 `json:"revenue"`
}
