package routes

import (
	"database/sql"
	handlers "madaurus/dev/assignment/app/handlers/files"
	"madaurus/dev/assignment/app/middlewares"

	"github.com/gin-gonic/gin"
)

func FilesRoute(c *gin.Engine, db *sql.DB) {
	FilesRoute := c.Group("/files", middlewares.Authentication())
	FilesRoute.GET(":fileid", handlers.GetFile())
}
