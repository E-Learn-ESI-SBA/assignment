package routes

import (
	"database/sql"
	handlers "madaurus/dev/assignment/app/handlers/assignments"
	"madaurus/dev/assignment/app/middlewares"

	"github.com/gin-gonic/gin"
)

func AssignmentsRoute(c *gin.Engine, db *sql.DB){
	assignmentsRoute := c.Group("/assignments", middlewares.Authentication())
	assignmentsRoute.GET("", handlers.GetAssignments(db))
	assignmentsRoute.GET(":assignmentId", handlers.GetAssignmentByID(db))
	assignmentsRoute.POST("", handlers.CreateAssignment(db))
	assignmentsRoute.PUT(":assignmentId", handlers.UpdateAssignment(db))
	assignmentsRoute.DELETE(":assignmentId", handlers.DeleteAssignmentByID(db))

	// assignmentsRoute.GET("/module/:moduleId", handlers.GetAssignmentsByModuleID(db))
	// assignmentsRoute.GET("/teacher/:teacherId", handlers.GetAssignmentsByTeacherID(db))
}