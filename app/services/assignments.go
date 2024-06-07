package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"madaurus/dev/assignment/app/interfaces"
	"madaurus/dev/assignment/app/models"
	"time"

	"github.com/google/uuid"
)

func CreateAssignment(ctx context.Context, db *sql.DB, assignment models.Assignment) error {
	_, err := db.ExecContext(ctx, "INSERT INTO assignments (id, title, description, deadline, year, module_id, teacher_id, file) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		assignment.ID, assignment.Title, assignment.Description, assignment.Deadline, assignment.Year, assignment.ModuleId, assignment.TeacherId, assignment.File)
	if err != nil {
		log.Printf("Error when creating assignments %v", err)
		return err
	}
	return nil
}

func GetAssignmentByID(ctx context.Context, db *sql.DB, assignmentId uuid.UUID) (*models.Assignment, error) {
	var assignment models.Assignment

	query := `SELECT id, title, description, file, deadline, year, teacher_id, module_id 
	          FROM assignments 
	          WHERE id = $1`
	err := db.QueryRowContext(ctx, query, assignmentId).Scan(&assignment.ID, &assignment.Title, &assignment.Description, &assignment.File, &assignment.Deadline, &assignment.Year, &assignment.TeacherId, &assignment.ModuleId)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No assignment with id %s", assignmentId)
			return nil, nil
		}
		log.Printf("Error getting assignment with id %s: %v", assignmentId, err)
		return nil, err
	}
	if assignment.Deadline.IsZero() {
        assignment.Deadline = time.Time{} 
    }
	return &assignment, nil
}

func UpdateAssignment(ctx context.Context, db *sql.DB, assignmentId uuid.UUID, editedAssignment models.Assignment) error {
	_, err := db.ExecContext(ctx, "UPDATE assignments SET title = $1, description = $2, deadline = $3, year = $4 WHERE id = $5",
		editedAssignment.Title, editedAssignment.Description, editedAssignment.Deadline, editedAssignment.Year, assignmentId)
	if err != nil {
		log.Printf("Error when updating assignment with ID %s: %v", assignmentId, err)
		return err
	}
	return nil
}

func DeleteAssignmentByID(ctx context.Context, db *sql.DB, assignmentID uuid.UUID) error {
	_, err := db.ExecContext(ctx, "DELETE FROM submissions WHERE assignment_id = $1", assignmentID)
	if err != nil {
	  return err 
	}
  
	_, err = db.ExecContext(ctx, "DELETE FROM assignments WHERE id = $1", assignmentID)
	if err != nil {
	  return err 
	}
	return nil
  }
  
func GetAssignments(ctx context.Context, db *sql.DB, filter interfaces.AssignmentFilter, filterId string, filterBy string) ([]models.Assignment, error) {
	var assignments []models.Assignment
	var query string
	var args []interface{}
	query = "SELECT id, title, description, file, deadline, year, teacher_id, module_id FROM assignments WHERE 1=1"
	// if *filter.ModuleId != "" {
	// 	query += " AND module_id = $1"
	// 	args = append(args, *filter.ModuleId)
	// }
	// if *filter.TeacherId != "" {
	// 	query += " AND teacher_id = $2"
	// 	args = append(args, *filter.TeacherId)
	// }
 
	if filterId != "" {
		if filterBy == "Year"{
			query += " AND year = $1"
			args = append(args, filterId)
		} else {
			query += " AND teacher_id = $1"
			args = append(args, filterId)
		}
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return assignments, err
	}

	defer rows.Close()
	for rows.Next() {
		var assignment models.Assignment
		if err := rows.Scan(&assignment.ID, &assignment.Title, &assignment.Description, &assignment.File, &assignment.Deadline, &assignment.Year, &assignment.TeacherId, &assignment.ModuleId); err != nil {
			log.Printf("Error scanning row: %v", err)
			return assignments, err
		}
		assignments = append(assignments, assignment)
		fmt.Print("\n")
		fmt.Print(assignment.File)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return assignments, err
	}

	return assignments, nil
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
