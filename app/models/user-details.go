package models

import "github.com/google/uuid"

type UserDetails struct {
	ID   uuid.UUID
	Name string
}
