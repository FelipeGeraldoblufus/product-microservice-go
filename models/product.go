package models

type Product struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	ProductID   string `gorm:"not null" json:"product_id"` 
	Name        string  `gorm:"not null;unique" json:"name"`
	Price       int `gorm:"not null" json:"price"`       // Cambia a float64 para representar precios
	Stock       int     `gorm:"not null" json:"stock"`       // Entero para la cantidad en stock
	Description string  `gorm:"not null" json:"description"`
	Category string `gorm:"not null" json:"category"` // Texto descriptivo del producto
}

type User struct {
	ID       uint       `gorm:"primaryKey" json:"id"`
	Username string     `gorm:"not null;unique" json:"username"`
}
