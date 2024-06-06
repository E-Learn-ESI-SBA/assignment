package services

import (
	"context"
	"database/sql"
	"log"
	"madaurus/dev/assignment/app/models"
	"time"

	"github.com/google/uuid"
)

func GetSubmissionByAssignmentID(ctx context.Context, db *sql.DB, assignmentId uuid.UUID, studentId string) ([]models.Submission, error) {
	var submissions []models.Submission
	var rows *sql.Rows
	var err error
	
	if studentId != ""{
		    rows, err = db.Query("SELECT * FROM submissions WHERE assignment_id = $1 AND student_id = $2", assignmentId, studentId)
	} else {
		rows, err = db.Query("SELECT * FROM submissions WHERE assignment_id = $1", assignmentId)
	}
	if err != nil {
		log.Printf("Error getting submissions with assignmentId %s: %v", assignmentId, err)
		return submissions, err
	}
	defer rows.Close()

	for rows.Next() {
		var submission models.Submission
		var feedback, evaluatedAt sql.NullString 

		if err := rows.Scan(
			&submission.ID,
			&submission.File,
			&submission.Grade,
			&feedback,
			&submission.AssignmentId,
			&submission.StudentId,
			&submission.CreatedAt,
			&submission.UpdatedAt,
			&evaluatedAt,
		); err != nil {
			log.Printf("Error scanning row: %v", err)
			return submissions, err
		}

		if feedback.Valid {
			submission.Feedback = feedback.String
		} else {
			submission.Feedback = "" 
		}

		if evaluatedAt.Valid {
			evaluatedTime, err := time.Parse(time.RFC3339, evaluatedAt.String)
			if err != nil {
				log.Printf("Error parsing EvaluatedAt: %v", err)
				return submissions, err
			}
			submission.EvaluatedAt = &evaluatedTime
		} else {
			submission.EvaluatedAt = &time.Time{} 
		}

		submissions = append(submissions, submission)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return submissions, err
	}

	return submissions, nil
}

func CreateSubmission(ctx context.Context, db *sql.DB, submission models.Submission) error {
	log.Println("assignment id", submission.AssignmentId)
	_, err := db.Exec("INSERT INTO submissions (file, student_id, assignment_id, created_at, updated_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
		submission.File, submission.StudentId, submission.AssignmentId)
	if err != nil {
		log.Printf("Error when creating submission: %v", err)
		return err
	}
	return nil
}

func UpdateSubmission(ctx context.Context, db *sql.DB, submissionID int, editedSubmission models.Submission) error {
	_, err := db.Exec("UPDATE submissions SET file = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND student_id = $3",
		editedSubmission.File, submissionID, editedSubmission.StudentId)
	if err != nil {
		log.Printf("Error when updating submission with ID %d: %v", submissionID, err)
		return err
	}
	return nil
}

func DeleteSubmissionByID(ctx context.Context, db *sql.DB, submissionID uuid.UUID) error {
	_, err := db.Exec("DELETE FROM submissions WHERE id = $1", submissionID)
	if err != nil {
		log.Printf("Error when deleting submission with ID %d: %v", submissionID, err)
		return err
	}
	return nil
}

func EvaluateSubmission(ctx context.Context, db *sql.DB, evaluatedSubmission models.Submission, submissionID uuid.UUID) error {
	_, err := db.Exec("UPDATE submissions SET grade = $1, feedback = $2 WHERE id = $3",
		evaluatedSubmission.Grade, evaluatedSubmission.Feedback, submissionID)
	if err != nil {
		log.Printf("Error when evaluating submission with ID %d: %v", submissionID, err)
		return err
	}
	return nil
}
func GetSubmissionByID(ctx context.Context, db *sql.DB, submissionID uuid.UUID) (*models.Submission, error) {
	var submission models.Submission
	var evaluatedAt sql.NullTime
	var feedback sql.NullString

	err := db.QueryRow("SELECT * FROM submissions WHERE id = $1", submissionID.String()).Scan(
		&submission.ID,
		&submission.File,
		&submission.Grade,
		&feedback, 
		&submission.AssignmentId,
		&submission.StudentId,
		&submission.CreatedAt,
		&submission.UpdatedAt,
		&evaluatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No submission with id %s", submissionID)
			return nil, nil
		}
		log.Printf("Error getting submission with id %s: %v", submissionID, err)
		return nil, err
	}

	if feedback.Valid {
		submission.Feedback = feedback.String
	}

	if evaluatedAt.Valid {
		submission.EvaluatedAt = &evaluatedAt.Time
	}

	return &submission, nil
}
