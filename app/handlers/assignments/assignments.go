package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"madaurus/dev/assignment/app/interfaces"
	"madaurus/dev/assignment/app/models"
	"madaurus/dev/assignment/app/services"
	"madaurus/dev/assignment/app/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAssignments(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var assignments []models.Assignment
		var teacherId int
		var err error

		if teacherIDStr := c.Query("teacher_id"); teacherIDStr != "" {
			teacherId, err = strconv.Atoi(teacherIDStr)
			if err != nil {
				c.JSON(400, gin.H{"error": "error when parsing teacherid"})
				return
			}
		}
		moduleId := c.Query("module_id")

		filter := interfaces.AssignmentFilter{
			Teacher: &teacherId,
			Module:  &moduleId,
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
		assignment.ID = user.ID
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

		assignmentId, err := strconv.Atoi(assignmentIdStr)
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

		assignmentId, errP := strconv.Atoi(assignmentIdStr)
		if errP != nil {
			c.JSON(400, gin.H{"error": errors.New("id is not valid")})
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
		assignmenIdStr, errr := c.Params.Get("assignmentId")
		fmt.Println(assignmenIdStr)
		if !errr {
			c.JSON(400, gin.H{"error": "Error when parsig id"})
			return
		}
		assignmentId, errP := strconv.Atoi(assignmenIdStr)
		if errP != nil {
			c.JSON(400, gin.H{"error": errors.New("id is not valid")})
			return
		}
		err := services.DeleteAssignmentByID(c.Request.Context(), db, assignmentId)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Assignment Deleted Successfully"})
	}
}
