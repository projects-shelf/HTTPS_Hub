package root

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func RootHandler(portMap map[string]string, domain string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Param("path")

		// api
		if path == "/api/domain" {
			c.JSON(http.StatusOK, domain)
			return
		}
		if path == "/api/ip" {
			c.JSON(http.StatusOK, os.Getenv("LOCAL_IP"))
			return
		}
		if path == "/api/config" {
			c.JSON(http.StatusOK, portMap)
			return
		}

		// static files
		if path == "" || path == "/" {
			path = "/index.html"
		}

		filePath := filepath.Join("/front", filepath.Clean(path))

		c.File(filePath)
	}
}
