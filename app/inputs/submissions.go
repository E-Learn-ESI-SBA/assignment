package inputs

import (
	"time"

	"github.com/google/uuid"
)

type NewSubmissionInput struct {
	File       string    `json:"file"`
	Assignment uuid.UUID       `json:"assignment_id"`
	Student    string       `json:"student_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}