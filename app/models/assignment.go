package models

import (
	"time"

	"github.com/google/uuid"
)

type Assignment struct {
	ID          uuid.UUID       `json:"id"`
	Title       string          `json:"title"`
	Description string    		`json:"description"`
	Files	    []string        `json:"files"`
	Deadline    time.Time 		`json:"deadline"`
	Promo       string       	`json:"promo"`
	Groups      []uuid.UUID     `json:"groups"`
	Teacher     uuid.UUID       `json:"teacher_id"`
	Module      uuid.UUID       `json:"module_id"`
}
