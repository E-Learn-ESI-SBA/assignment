package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"github.com/gin-gonic/gin"
)

func GetFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		fileID := c.Param("fileid")
		// toAdd: uuid validation


		filePath := filepath.Join("uploads", fileID)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		c.File(filePath)
	}
}
