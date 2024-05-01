package services

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Submission struct {
	ID           int     	`json:"id"`
	File         string  	`json:"file"`
	DateTime	 time.Time 	`json:"date_time"`
	Grade        float64 	`json:"grade"`
	Feedback     string  	`json:"feedback"`
	StudentID    int     	`json:"student_id"`
	AssignmentID int     	`json:"assignment_id"`
}

func CreateSubmission(ctx context.Context, db *sql.DB, submission Submission) error {
	_, err := db.Exec("INSERT INTO submissions (file, grade, feedback, student_id, assignment_id) VALUES ($1, $2, $3, $4, $5)",
		submission.File, submission.Grade, submission.Feedback, submission.StudentID, submission.AssignmentID)
	if err != nil {
		log.Printf("Error when creating submission: %v", err)
		return err
	}
	return nil
}

func UpdateSubmission(ctx context.Context, db *sql.DB, submissionID int, editedSubmission Submission) error {
	_, err := db.Exec("UPDATE submissions SET file = $1, grade = $2, feedback = $3, student_id = $4, assignment_id = $5 WHERE id = $6",
		editedSubmission.File, editedSubmission.Grade, editedSubmission.Feedback, editedSubmission.StudentID, editedSubmission.AssignmentID, submissionID)
	if err != nil {
		log.Printf("Error when updating submission with ID %d: %v", submissionID, err)
		return err
	}
	return nil
}

func DeleteSubmission(ctx context.Context, db *sql.DB, submissionID int) error {
	_, err := db.Exec("DELETE FROM submissions WHERE id = $1", submissionID)
	if err != nil {
		log.Printf("Error when deleting submission with ID %d: %v", submissionID, err)
		return err
	}
	return nil
}
