package models

import (
	"time"
)

type Assignment struct {
	ID          int       `json:"id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Promo       int       `json:"promo" binding:"required"`
	Groups      []int     `json:"groups" binding:"required"`
	Teacher     int       `json:"teacher_id" binding:"required"`
	Module      string    `json:"module_id" binding:"required"`
}
