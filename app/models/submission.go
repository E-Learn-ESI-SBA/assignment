package models

import (
	"time"

	"github.com/google/uuid"
)

type Submission struct {
	ID           uuid.UUID        	`json:"id"`
	File         string    		 	`json:"file"`
	Grade        float32   			`json:"grade"`
	Feedback     string    			`json:"feedback"`
    Assignment   uuid.UUID       	`json:"assignment_id"`
	Student      uuid.UUID       	`json:"student_id"`
	CreatedAt    time.Time 			`json:"created_at"`
	UpdatedAt    time.Time 			`json:"updated_at"`
	EvaluatedAt  time.Time 			`json:"evaluated_at"`
}
