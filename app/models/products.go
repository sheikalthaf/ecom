package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID      `gorm:"primaryKey;type:uuid;column:Id" json:"Id"`
	Name        string         `gorm:"type:varchar(200);column:Name" json:"Name"`
	CategoryID  uuid.UUID      `gorm:"type:uuid;column:CategoryId" json:"CategoryId"`
	Category    Category       `gorm:"foreignKey:CategoryId" json:"Category"`
	ProductCode string         `gorm:"type:varchar(200);unique;column:ProductCode" json:"ProductCode"`
	Price       float64        `gorm:"type:decimal;column:Price" json:"Price"`
	Description string         `gorm:"column:Description" json:"Description"`
	Size        string         `gorm:"type:varchar(200);column:Size" json:"Size"`
	Quantity    int            `gorm:"type:int;column:Quantity" json:"Quantity"`
	IsDeleted   bool           `gorm:"type:bool;column:IsDeleted" json:"IsDeleted"`
	Image       pq.StringArray `gorm:"type:varchar(200)[];column:Image" json:"Image"`
	CreatedAt   time.Time      `gorm:"column:CreatedAt" json:"CreatedAt"`
	Tags        pq.StringArray `gorm:"type:varchar(200)[];column:Tags" json:"Tags"`
}

func (product *Product) TableName() string {
	return "Products"
}

var ProductColumns = struct {
	ID          string
	Name        string
	ProductCode string
	Price       string
	Size        string
	Quantity    string
	IsDeleted   string
	Image       string
	CreatedAt   string
	Tags        string
}{
	ID:          "Id",
	Name:        "Name",
	ProductCode: "ProductCode",
	Price:       "Price",
	Size:        "Size",
	Quantity:    "Quantity",
	IsDeleted:   "IsDeleted",
	Image:       "Image",
	CreatedAt:   "CreatedAt",
	Tags:        "Tags",
}

func (product *Product) BeforeCreate(tx *gorm.DB) (err error) {
	product.ID = uuid.New()
	return
}
