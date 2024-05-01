package services

import (
	"context"
	"database/sql"
	"log"
	"madaurus/dev/assignment/app/inputs"
	"madaurus/dev/assignment/app/models"
)

func GetSubmissionByAssignmentID(ctx context.Context, db *sql.DB, assignmentId int) ([]models.Submission, error) {
	var submissions []models.Submission

	rows, err := db.Query("SELECT * FROM submissions WHERE assignment_id = $1", assignmentId)
	if err != nil {
		log.Printf("Error geting submissions with assignmentId %d", assignmentId)
		return submissions, err
	}
	for rows.Next() {
		var submission models.Submission
		if err = rows.Scan(&submission.ID, &submission.Assignment, &submission.Feedback, &submission.File, &submission.Grade, &submission.Student, &submission.CreatedAt , &submission.UpdatedAt, &submission.EvaluatedAt); err != nil {
			log.Printf("Error scanning row: %v", err)
			return submissions, err
		}
		submissions = append(submissions, submission)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return submissions, err
	}

	return submissions, nil
	
}


func CreateSubmission(ctx context.Context, db *sql.DB, submission inputs.NewSubmissionInput) error {
	_, err := db.Exec("INSERT INTO submissions (file, student_id, assignment_id, created_at, updated_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
		submission.File, submission.Student, submission.Assignment)
	if err != nil {
		log.Printf("Error when creating submission: %v", err)
		return err
	}
	return nil
}

func UpdateSubmission(ctx context.Context, db *sql.DB, submissionID int, editedSubmission models.Submission) error {
	_, err := db.Exec("UPDATE submissions SET file = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND student_id = $3",
		editedSubmission.File, submissionID, editedSubmission.Student)
	if err != nil {
		log.Printf("Error when updating submission with ID %d: %v", submissionID, err)
		return err
	}
	return nil
}

func DeleteSubmissionByID(ctx context.Context, db *sql.DB, submissionID int) error {
	_, err := db.Exec("DELETE FROM submissions WHERE id = $1", submissionID)
	if err != nil {
		log.Printf("Error when deleting submission with ID %d: %v", submissionID, err)
		return err
	}
	return nil
}


func EvaluateSubmission(ctx context.Context, db *sql.DB,evaluatedSubmission models.Submission,submissionID int) error {
	_, err := db.Exec("UPDATE submissions SET grade = $1, feedback = $2 WHERE id = $3",
	evaluatedSubmission.Grade, evaluatedSubmission.Feedback, submissionID)
	if err != nil {
		log.Printf("Error when evaluating submission with ID %d: %v", submissionID, err)
		return err
	}
	return nil
}



func GetSubmissionByID(ctx context.Context, db *sql.DB, submissionId int) (models.Submission, error) {
	var submission models.Submission

	err := db.QueryRow("SELECT * FROM submissions WHERE id = $1", submissionId).Scan(submission.ID, submission.File,
		submission.Grade, submission.Feedback, submission.Student, submission.CreatedAt,
		submission.UpdatedAt, submission.EvaluatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No submission with id %d", submissionId)
			return submission, nil
		}
		log.Printf("Error geting submission with id %d", submissionId)
		return submission, err
	}
	return submission, nil
}