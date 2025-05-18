package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func ProxyHandler(portMap map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Host
		key := host

		if idx := strings.Index(host, "."); idx != -1 {
			key = host[:idx]
		} else {
			c.String(404, "No key found")
			return
		}

		port, ok := portMap[key]
		if !ok {
			c.String(404, "No port mapping found for key: %s", key)
			return
		}

		domain := port
		if strings.HasPrefix(port, ":") {
			ip := os.Getenv("LOCAL_IP") // e.g., 192.168.0.2
			domain = ip + port
		}

		path := c.Param("path")
		target_url := "http://" + domain + path
		query := c.Request.URL.RawQuery
		if query != "" {
			target_url += "?" + query
		}

		target, err := url.Parse(target_url)
		if err != nil {
			c.String(http.StatusInternalServerError, "Invalid target URL: %s", err.Error())
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(target)

		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.URL.Path = c.Request.URL.Path
			req.URL.RawQuery = c.Request.URL.RawQuery
			req.Header = c.Request.Header.Clone()
			req.Host = target.Host
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
