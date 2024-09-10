package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;column:Id" json:"Id"`
	Name      string    `gorm:"type:varchar(255);column:Name" json:"Name"`
	MobileNo  string    `gorm:"type:varchar(255);column:MobileNo" json:"MobileNo"`
	Email     string    `gorm:"type:varchar(255);column:Email" json:"Email"`
	Password  string    `gorm:"type:varchar(255);column:Password" json:"Password"`
	Role      string    `gorm:"type:varchar(255);column:Role" json:"Role"`
	IsDeleted bool      `gorm:"type:boolean;column:IsDeleted" json:"IsDeleted"`
	CreatedAt time.Time `gorm:"type:timestamp;column:CreatedAt" json:"CreatedAt"`
	UpdatedAt time.Time `gorm:"type:timestamp;column:UpdatedAt" json:"UpdatedAt"`
}

func (user *User) TableName() string {
	return "Users"
}

var UserColumns = struct {
	ID        string
	Name      string
	MobileNo  string
	Email     string
	Password  string
	Role      string
	IsDeleted string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "Id",
	Name:      "Name",
	MobileNo:  "MobileNo",
	Email:     "Email",
	Password:  "Password",
	Role:      "Role",
	IsDeleted: "IsDeleted",
	CreatedAt: "CreatedAt",
	UpdatedAt: "UpdatedAt",
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	user.IsDeleted = false
	user.Role = "user"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return
}
