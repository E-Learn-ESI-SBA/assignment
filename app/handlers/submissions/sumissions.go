package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"log"

	// "log"
	"madaurus/dev/assignment/app/models"
	"madaurus/dev/assignment/app/services"
	"madaurus/dev/assignment/app/shared"
	"madaurus/dev/assignment/app/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateSubmission(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.New().String() + ".zip"
		var filePath string

		user := c.MustGet("user").(*utils.UserDetails)
		fmt.Println(user.ID)
		assignmentIDStr := c.Param("assignmentId")

		assignmentID, err := uuid.Parse(assignmentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
			return
		}

		err = c.Request.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil && err != http.ErrNotMultipart {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
			return
		}

		file, _, err := c.Request.FormFile("file")
		if err != nil {
			filePath = ""
		} else {
			defer file.Close()

			// save the file to fs
			fileDir := "./uploads"
			if _, err := os.Stat(fileDir); os.IsNotExist(err) {
				err = os.Mkdir(fileDir, os.ModePerm)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file directory"})
					return
				}
			}
			filePath = filepath.Join(fileDir, id)
			out, err := os.Create(filePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
				return
			}
			defer out.Close()

			_, err = io.Copy(out, file)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
				return
			}
		}

		var submission models.Submission
		submission.StudentId = user.ID
		submission.AssignmentId = assignmentID
		submission.CreatedAt = time.Now()
		submission.UpdatedAt = time.Now()
		submission.File = id // if no file uploaded -> File = ""
		// err = c.ShouldBindJSON(&submission)
		// if err != nil {
		// 	log.Println(err)
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse submission data"})
		// 	return
		// }

		err = services.CreateSubmission(c.Request.Context(), db, submission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create submission"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Submission Created Successfully"})
	}
}

func UpdateSubmission(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var editedSubmission models.Submission
		err := c.BindJSON(&editedSubmission)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		submissionIdStr, exists := c.Params.Get("submissionId")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing submission ID"})
			return
		}

		submissionId, err := strconv.Atoi(submissionIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID format"})
			return
		}

		err = services.UpdateSubmission(c.Request.Context(), db, submissionId, editedSubmission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update submission"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Submission Updated Successfully"})
	}
}

func GetSubmissionsByAssignmentId(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var submissions []models.Submission
		
		value, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"message": shared.UNAUTHORIZED})
			return
		}
		
		user, ok := value.(*utils.UserDetails)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse user details"})
			return
		}
		assignmentIDStr, exists := c.Params.Get("assignmentId")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing assignment ID"})
			return
		}

		assignmentId, err := uuid.Parse(assignmentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID format"})
			return
		}

		var userID =""
		if user.Role == "student" {
			userID = user.ID
		}
		submissions, err = services.GetSubmissionByAssignmentID(c.Request.Context(), db, assignmentId, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve submissions"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"submissions": submissions})
	}
}

func GetSubmissionByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		submissionIdStr, exists := c.Params.Get("submissionId")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing submission ID"})
			return
		}

		submissionId, err := uuid.Parse(submissionIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID format"})
			return
		}
		submission, err := services.GetSubmissionByID(c.Request.Context(), db, submissionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve submission"})
			return
		}
		if submission == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"submission": submission})
	}
}

func DeleteSubmissionByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		submissionIdStr, exists := c.Params.Get("submissionId")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing submission ID"})
			return
		}

		submissionId, err := uuid.Parse(submissionIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID format"})
			return
		}

		err = services.DeleteSubmissionByID(c.Request.Context(), db, submissionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete submission"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Submission Deleted Successfully"})
	}
}

func EvaluateSubmission(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var evaluatedSubmission models.Submission
		submissionIDStr := c.Param("submissionId")

		submissionId, err := uuid.Parse(submissionIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID format"})
			return
		}

		if err := c.BindJSON(&evaluatedSubmission); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse submission data"})
			return
		}

		err = services.EvaluateSubmission(c.Request.Context(), db, evaluatedSubmission, submissionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to evaluate submission"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Submission Evaluated Successfully"})
	}
}
