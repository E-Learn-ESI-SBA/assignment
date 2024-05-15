package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"madaurus/dev/assignment/app/inputs"
	"madaurus/dev/assignment/app/models"
	"madaurus/dev/assignment/app/services"
	"madaurus/dev/assignment/app/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetSubmissionsByAssignmentId(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var submissions []models.Submission
		var err error

		assignmentIDStr, errP := c.Params.Get("assignmentId")
		if !errP {
			c.JSON(400, gin.H{"error": "Error when parsing assignmentId"})
			return
		}

		assignmentId, err := uuid.Parse(assignmentIDStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Error when parsing assignmentId"})
			return
		}

		submissions, err = services.GetSubmissionByAssignmentID(c.Request.Context(), db, assignmentId)
		if err != nil {
			c.JSON(400, gin.H{"error": "Something Went Wrong"})
			return
		}
		c.JSON(200, gin.H{"message": submissions})

	}
}

func CreateSubmission(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var submission inputs.NewSubmissionInput
		user := c.MustGet("user").(*utils.UserDetails)
		assignmentIDStr := c.Param("assignmentId")

		assignmentID, err := uuid.Parse(assignmentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
			return
		}

		if err := c.BindJSON(&submission); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse submission data"})
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file from request"})
			return
		}

		dst := "uploads/submissions/" + file.Filename
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		submission.Student = user.ID
		submission.Assignment = assignmentID
		submission.CreatedAt = time.Now()
		submission.UpdatedAt = time.Now()
		submission.File = dst 

		err = services.CreateSubmission(c.Request.Context(), db, submission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create submission"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Submission Created Successfully"})
	}
}

func UpdateSubmission(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var editedSubmission models.Submission
		err := c.BindJSON(&editedSubmission)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		submissionIdStr, errP := c.Params.Get("submissionId")

		if !errP {
			c.JSON(400, gin.H{"error": "error when parsing submission id"})
			return
		}

		submissionId, err := strconv.Atoi(submissionIdStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "error when parsing submission id"})
			return
		}

		err = services.UpdateSubmission(c.Request.Context(), db, submissionId, editedSubmission)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Submission Updated Successfully"})
	}
}

func GetSubmissionByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		submissionIdStr, errr := c.Params.Get("submissionId")
		fmt.Println(submissionIdStr)
		if !errr {
			c.JSON(400, gin.H{"error": "Error when parsig id"})
			return
		}

		submissionId, errP := strconv.Atoi(submissionIdStr)
		if errP != nil {
			c.JSON(400, gin.H{"error": errors.New("id is not valid")})
			return
		}
		submission, err := services.GetSubmissionByID(c.Request.Context(), db, submissionId)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": submission})
	}
}

func DeleteSubmissionByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		submissionIdStr, errr := c.Params.Get("submissionId")
		fmt.Println(submissionIdStr)
		if !errr {
			c.JSON(400, gin.H{"error": "Error when parsig id"})
			return
		}
		submissionId, errP := strconv.Atoi(submissionIdStr)
		if errP != nil {
			c.JSON(400, gin.H{"error": errors.New("id is not valid")})
			return
		}
		err := services.DeleteSubmissionByID(c.Request.Context(), db, submissionId)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Submission Deleted Successfully"})
	}
}
