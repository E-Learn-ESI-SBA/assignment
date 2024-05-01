package inputs

import "time"

type NewSubmissionInput struct {
	File       string    `json:"file"`
	Assignment int       `json:"assignment_id"`
	Student    int       `json:"student_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}