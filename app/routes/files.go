package routes

import (
	"database/sql"
	handlers "madaurus/dev/assignment/app/handlers/files"

	"github.com/gin-gonic/gin"
)

func FilesRoute(c *gin.Engine, db *sql.DB) {
	FilesRoute := c.Group("/files")
	FilesRoute.GET(":fileid", handlers.GetFile())
}
