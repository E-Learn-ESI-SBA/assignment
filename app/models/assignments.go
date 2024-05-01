package models

import (
	"time"
)

type Assignment struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Promo       int       `json:"promo"`
	Groups      []int     `json:"groups"`
	Teacher     int       `json:"teacher_id"`
	Module      string    `json:"module_id"`
}
