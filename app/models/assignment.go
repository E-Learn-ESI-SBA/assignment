package models

import (
	"time"

	"github.com/google/uuid"
)

type Assignment struct {
	ID          uuid.UUID       `json:"id"`
	Title       string          `json:"title"`
	Description string    		`json:"description"`
	File	    string          `json:"file"`
	Deadline    time.Time 		`json:"deadline"`
	Year        string       	`json:"year"`
	TeacherId   string          `json:"teacher_id"`
	ModuleId    string          `json:"module_id"`
}
