package models

import "github.com/google/uuid"

type Model struct {
	Id uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
}
