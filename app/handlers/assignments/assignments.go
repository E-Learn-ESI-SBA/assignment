package handlers

import (
	"database/sql"
	"fmt"
	"madaurus/dev/assignment/app/interfaces"
	"madaurus/dev/assignment/app/models"
	"madaurus/dev/assignment/app/services"
	"madaurus/dev/assignment/app/utils"
	"net/http"

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
	return func(c *gin.Context) {
		var assignment models.Assignment
		user := c.MustGet("user").(*utils.UserDetails)
		err := c.BindJSON(&assignment)

		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err = services.CreateAssignment(c.Request.Context(), db, assignment, user.ID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Assignment Created Successfully"})
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
		assignmentIdStr, errP := c.Params.Get("assignmentId")
		if !errP {
			c.JSON(400, gin.H{"error": "error when parsing assignment id"})
			return
		}

		assignmentId, err := uuid.Parse(assignmentIdStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "error when parsing assignment id"})
			return
		}
		err = services.UpdateAssignment(c.Request.Context(), db, assignmentId, editedAssignment)
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
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": assignment})
	}
}


func AddAssignmentFile(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file from request"})
			return
		}

		dst := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		assignmentIDStr := c.Param("assignmentId")
		assignmentID, err := uuid.Parse(assignmentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
			return
		}

		err = services.AddAssignmentFile(c.Request.Context(), db, assignmentID, dst)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update database with file link"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "File uploaded and assignment updated successfully"})
	}
}

func DeleteAssignmentFile(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		assignmentIDStr := c.Param("assignmentId")
		fileId := c.Param("fileId")
		assignmentID, err := uuid.Parse(assignmentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
			return
		}

		err = services.DeleteAssignmentFile(c.Request.Context(), db, fileId, assignmentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete assignment file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
	}
}

func GetAssignmentFile(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		basePath := "uploads/"       
		fileID := c.Param("fileId") 

		filePath := basePath + fileID
		c.File(filePath)
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
