package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;column:Id" json:"Id"`
	UserID    uuid.UUID `gorm:"type:uuid;column:UserId" json:"UserId"`
	ProductID uuid.UUID `gorm:"type:uuid;column:ProductId" json:"ProductId"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"Product"`
	Quantity  int       `gorm:"type:int;column:Quantity" json:"Quantity"`
	CreatedAt time.Time `gorm:"column:CreatedAt" json:"CreatedAt"`
}

func (cart *Cart) TableName() string {
	return "Carts"
}

var CartColumns = struct {
	ID        string
	UserID    string
	ProductID string
	Quantity  string
	CreatedAt string
}{
	ID:        "Id",
	UserID:    "UserId",
	ProductID: "ProductId",
	Quantity:  "Quantity",
	CreatedAt: "CreatedAt",
}

func (cart *Cart) BeforeCreate(tx *gorm.DB) (err error) {
	cart.ID = uuid.New()
	return
}
