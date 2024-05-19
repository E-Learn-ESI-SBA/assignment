package routes

import (
	"database/sql"
	handlers "madaurus/dev/assignment/app/handlers/submissions"
	"madaurus/dev/assignment/app/middlewares"

	"github.com/gin-gonic/gin"
)

func SubmissionsRoute(c *gin.Engine, db *sql.DB) {
	SubmissionsRoute := c.Group("/assignments/:assignmentId/submissions/", middlewares.Authentication())
	SubmissionsRoute.GET("", handlers.GetSubmissionsByAssignmentId(db))
	SubmissionsRoute.GET(":submissionId", handlers.GetSubmissionByID(db))
	SubmissionsRoute.POST("", handlers.CreateSubmission(db))
	SubmissionsRoute.PUT(":submissionId", handlers.EvaluateSubmission(db))
	// SubmissionsRoute.PUT(":submissionId", handlers.UpdateSubmission(db))
	SubmissionsRoute.DELETE(":submissionId", handlers.DeleteSubmissionByID(db))
}
