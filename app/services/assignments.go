package services

import (
	"context"
	"database/sql"
	"log"
	"madaurus/dev/assignment/app/interfaces"
	"madaurus/dev/assignment/app/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
) 

func CreateAssignment(ctx context.Context, db *sql.DB, assignment models.Assignment, teacherID string) error {
	log.Printf("module id %v", assignment.ModuleId)
	_, err := db.Exec("INSERT INTO assignments (title,description,deadline,year,groups,module_id,teacher_id,files) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
		assignment.Title, assignment.Description, assignment.Deadline, assignment.Year, pq.Array(assignment.Groups), assignment.ModuleId, teacherID, pq.Array(assignment.Files))
	if err != nil {
		log.Printf("Error when creating assignments %v", err)
		return err
	}
	return nil
}
 
func GetAssignmentByID(ctx context.Context, db *sql.DB, assignmentId uuid.UUID) (models.Assignment, error) {
	var assignment models.Assignment

	err := db.QueryRow("SELECT * FROM assignments WHERE id = $1", assignmentId).Scan(assignment.ID, assignment.Title,
		assignment.Description, assignment.Deadline,
		assignment.Year, assignment.Groups, assignment.TeacherId,
		assignment.ModuleId)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No assignment with id %d", assignmentId)
			return assignment, nil
		}
		log.Printf("Error geting assignment with id %d", assignmentId)
		return assignment, err
	}
	return assignment, nil
}

func UpdateAssignment(ctx context.Context, db *sql.DB, assignmentId uuid.UUID, editedAssignment models.Assignment) error {
	_, err := db.Exec("UPDATE assignments SET title = $1, description = $2, deadline = $3, year = $4, groups = $5 WHERE id = $6",
		editedAssignment.Title, editedAssignment.Description, editedAssignment.Deadline, editedAssignment.Year, pq.Array(editedAssignment.Groups), assignmentId)
	if err != nil {
		log.Printf("Error when updating assignment with ID %d: %v", assignmentId, err)
		return err
	}
	return nil
}

func DeleteAssignmentByID(ctx context.Context, db *sql.DB, assignmentID uuid.UUID) error {
	_, err := db.Exec("DELETE FROM assignments WHERE id = $1", assignmentID)
	if err != nil {
		log.Printf("Error when deleting assignment with ID %d: %v", assignmentID, err)
		return err
	}
	return nil
}

func GetAssignments(ctx context.Context, db *sql.DB, filter interfaces.AssignmentFilter) ([]models.Assignment, error) {
	var assignments []models.Assignment
	var query string
	var args []interface{}

	query = "SELECT * FROM assignments WHERE 1=1"
	if filter.ModuleId != nil {
		query += " AND module_id = $1"
		args = append(args, filter.ModuleId)
	}
	if filter.TeacherId != nil {
		query += " AND teacher_id = $2"
		args = append(args, filter.TeacherId)
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("Error escuting query %v", err)
		return assignments, err
	}

	defer rows.Close()

	for rows.Next() {
		var assignment models.Assignment
		if err := rows.Scan(assignment.ID, assignment.Title, assignment.Description, assignment.Deadline, assignment.Year, assignment.Groups, assignment.TeacherId, assignment.ModuleId); err != nil {
			log.Printf("Error scanning row: %v", err)
			return assignments, err
		}
		assignments = append(assignments, assignment)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return assignments, err
	}

	return assignments, nil
}



func AddAssignmentFile(ctx context.Context, db *sql.DB, assignmentId uuid.UUID, linkFile string) (error) {
	_, err := db.Exec("UPDATE assignments SET files = array_append(files, $1) WHERE id = $2", linkFile, assignmentId)
	if err != nil {
		log.Printf("Error updating assignments table with link file: %v", err)
		return  err
	}
	return nil
}

func DeleteAssignmentFile(ctx context.Context, db *sql.DB, fileId string, assignmentId uuid.UUID) error {
	_, err := db.Exec("UPDATE assignments SET files = array_remove(files, $1) WHERE id = $2", fileId, assignmentId)
	if err != nil {
		log.Printf("Error removing file link from assignments table: %v", err)
		return err
	}
	return nil
}

// func GetAssignmentsByTeacherID(ctx context.Context, db *sql.DB, teacherId int) ([]models.Assignment, error) {
// 	var assignments []models.Assignment
// 	rows, err := db.Query("SELECT * FROM assignments WHERE id = $1", teacherId)
// 	if err != nil {
// 		log.Printf("There is no assignment with teacher_id %d", teacherId)
// 		return assignments, err
// 	}
// 	for rows.Next() {
// 		var assignment models.Assignment
// 		if err := rows.Scan(assignment.ID, assignment.Title,
// 			assignment.Description, assignment.Deadline,
// 			assignment.Year, assignment.Groups, assignment.Teacher,
// 			assignment.Module); err != nil {
// 			log.Printf("Error %v", err)
// 			return assignments, err
// 		}
// 		assignments = append(assignments, assignment)
// 	}
// 	if err := rows.Err(); err != nil {
// 		log.Printf("Error %v", err)
// 		return assignments, err
// 	}
// 	return assignments, nil
// }

// func GetAssignmentsByModule(ctx context.Context, db *sql.DB, moduleId int) ([]models.Assignment, error) {
// 	var assignments []models.Assignment
// 	rows, err := db.Query("SELECT * FROM assignments WHERE module_id = $1", moduleId)
// 	if err != nil {
// 		log.Printf("There is no assignment with module_id %d", moduleId)
// 		return assignments, err
// 	}
// 	for rows.Next() {
// 		var assignment models.Assignment
// 		if err := rows.Scan(assignment.ID, assignment.Title,
// 			assignment.Description, assignment.Deadline,
// 			assignment.Year, assignment.Groups, assignment.Teacher,
// 			assignment.Module); err != nil {
// 			log.Printf("Error geting assignment with module_id= %d", moduleId)
// 			return assignments, err
// 		}
// 		assignments = append(assignments, assignment)
// 	}
// 	return assignments, err
// }
