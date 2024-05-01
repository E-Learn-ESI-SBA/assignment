package models

import (
	"time"
)

type Submission struct {
	ID           int       `json:"id"`
	File         string    `json:"file"`
	Grade        float32   `json:"grade"`
	Feedback     string    `json:"feedback"`
	Assignment   int       `json:"assignment_id"`
	Student      int       `json:"student_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	EvaluatedAt  time.Time `json:"evaluated_at"`
}
