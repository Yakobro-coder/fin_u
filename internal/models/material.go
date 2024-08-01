package models

import (
	"time"

	"github.com/google/uuid"
)

type Material struct {
	UUID          uuid.UUID `json:"uuid"`
	TypeMaterials string    `json:"type_materials"`
	Status        string    `json:"status"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
