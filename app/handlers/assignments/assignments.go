package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"madaurus/dev/assignment/app/interfaces"
	"madaurus/dev/assignment/app/models"
	"madaurus/dev/assignment/app/services"
	"madaurus/dev/assignment/app/shared"
	"madaurus/dev/assignment/app/utils"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAssignments(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var assignments []models.Assignment
		var err error

		teacherId := c.Query("teacher_id")
		moduleId := c.Query("module_id")

		filter := interfaces.AssignmentFilter{
			TeacherId: &teacherId,
			ModuleId:  &moduleId,
		}

		assignments, err = services.GetAssignments(c.Request.Context(), db, filter)
		if err != nil {
			c.JSON(400, gin.H{"error": "error, something went wrong"})
			return
		}
		c.JSON(200, gin.H{"message": assignments})

	}
}

func CreateAssignment(db *sql.DB) gin.HandlerFunc {
	return func(g *gin.Context) {
		id := uuid.New().String()
		value, exists := g.Get("user")
		if !exists {
			g.JSON(http.StatusInternalServerError, gin.H{"message": shared.UNAUTHORIZED})
			return
		}
		user := value.(*utils.UserDetails)

		err := g.Request.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil && err != http.ErrNotMultipart {
			g.JSON(http.StatusBadRequest, gin.H{"message": shared.FILE_TOO_LARGE})
			return
		}

		var filePath string
		file, _, err := g.Request.FormFile("file")
		if err != nil {
			filePath = ""
		} else {
			defer file.Close()

			// save the file to fs
			fileDir := "./uploads"
			if _, err := os.Stat(fileDir); os.IsNotExist(err) {
				err = os.Mkdir(fileDir, os.ModePerm)
				if err != nil {
					g.JSON(http.StatusInternalServerError, gin.H{"message": shared.FILE_NOT_CREATED})
					return
				}
			}
			filePath = filepath.Join(fileDir, id)
			out, err := os.Create(filePath)
			if err != nil {
				g.JSON(http.StatusInternalServerError, gin.H{"message": shared.FILE_NOT_CREATED})
				return
			}
			defer out.Close()

			_, err = io.Copy(out, file)
			if err != nil {
				g.JSON(http.StatusInternalServerError, gin.H{"message": shared.FILE_NOT_CREATED})
				return
			}
		}

		var assignment models.Assignment
		err = g.ShouldBind(&assignment)
		if err != nil {
			log.Printf("Error binding assignment: %v", err.Error())
			g.JSON(http.StatusNotAcceptable, gin.H{"message": shared.UNABLE_TO_PARSE})
			return
		}

		assignment.ID = uuid.New()
		assignment.File = id // if no file uploaded -> file = empty
		assignment.TeacherId = user.ID

		err = services.CreateAssignment(g.Request.Context(), db, assignment)
		if err != nil {
			log.Printf("Error creating assignment: %v", err.Error())
			g.JSON(http.StatusBadRequest, gin.H{"message": "Assignment not created"})
			return
		}

		g.JSON(http.StatusCreated, gin.H{"message": "Assignment Created Successfully"})
	}
}

func UpdateAssignment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var editedAssignment models.Assignment
		err := c.BindJSON(&editedAssignment)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		assignmentIDStr, exists := c.Params.Get("assignmentId")
		if !exists {
			c.JSON(400, gin.H{"error": "error when parsing assignment ID"})
			return
		}

		assignmentID, err := uuid.Parse(assignmentIDStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "error when parsing assignment ID"})
			return
		}

		err = services.UpdateAssignment(c.Request.Context(), db, assignmentID, editedAssignment)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Assignment Updated Successfully"})
	}
}

func GetAssignmentByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		assignmentIdStr, errr := c.Params.Get("assignmentId")
		fmt.Println(assignmentIdStr)
		if !errr {
			c.JSON(400, gin.H{"error": "Error when parsig id"})
			return
		}

		assignmentId, err := uuid.Parse(assignmentIdStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "error when parsing assignment id"})
			return
		}
		assignment, err := services.GetAssignmentByID(c.Request.Context(), db, assignmentId)
		if assignment == nil {
			c.JSON(404, gin.H{"error": shared.ASSIGNMENT_NOT_FOUND})
			return
		}
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": assignment})
	}
}

// func GetAssignmentsByTeacherID(db *sql.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		teacherIdStr, errr := c.Params.Get("teacherId")
// 		fmt.Println(teacherIdStr)
// 		if !errr {
// 			c.JSON(400, gin.H{"error": "Error when parsig id"})
// 			return
// 		}
// 		teacherId, errP := strconv.Atoi(teacherIdStr)
// 		if errP != nil {
// 			c.JSON(400, gin.H{"error": errors.New("id is not valid")})
// 			return
// 		}
// 		assignments, err := services.GetAssignmentsByTeacherID(c.Request.Context(), db, teacherId)
// 		if err != nil {
// 			c.JSON(400, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(200, gin.H{"message": assignments})
// 	}
// }

// func GetAssignmentsByModuleID(db *sql.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		moduleIdStr, errr := c.Params.Get("moduleId")
// 		fmt.Println(moduleIdStr)
// 		if !errr {
// 			c.JSON(400, gin.H{"error": "Error when parsig id"})
// 			return
// 		}
// 		moduleId, errP := strconv.Atoi(moduleIdStr)
// 		if errP != nil {
// 			c.JSON(400, gin.H{"error": errors.New("id is not valid")})
// 			return
// 		}
// 		assignments, err := services.GetAssignmentsByTeacherID(c.Request.Context(), db, moduleId)
// 		if err != nil {
// 			c.JSON(400, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(200, gin.H{"message": assignments})
// 	}
// }

func DeleteAssignmentByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		assignmentIdStr, errr := c.Params.Get("assignmentId")
		if !errr {
			c.JSON(400, gin.H{"error": "Error when parsig id"})
			return
		}

		assignmentId, err := uuid.Parse(assignmentIdStr)

		if err != nil {
			c.JSON(400, gin.H{"error": "error when parsing assignment id"})
			return
		}
		err = services.DeleteAssignmentByID(c.Request.Context(), db, assignmentId)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Assignment Deleted Successfully"})
	}
}
