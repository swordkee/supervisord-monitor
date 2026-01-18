package main

import (
	"crypto/subtle"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"supervisord-monitor/config"
	"supervisord-monitor/handlers"

	"github.com/gin-gonic/gin"
)

func getMIMEType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".js":
		return "text/javascript; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".html":
		return "text/html; charset=utf-8"
	case ".mp3":
		return "audio/mpeg"
	case ".png", ".jpg", ".jpeg", ".gif", ".ico", ".webp":
		return "image/" + ext[1:]
	case ".woff", ".woff2":
		return "font/" + ext[1:]
	case ".json":
		return "application/json"
	case ".txt":
		return "text/plain; charset=utf-8"
	default:
		return mime.TypeByExtension(ext)
	}
}

func main() {
	configPath := flag.String("config", "", "Path to config file")
	port := flag.Int("port", 0, "Port to listen on (overrides config)")
	flag.Parse()

	if err := config.LoadConfig(*configPath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	listenPort := config.Cfg.Port
	if *port > 0 {
		listenPort = *port
	}

	subFS, err := fs.Sub(FrontendFS, "frontend/dist")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(handlers.LoggerMiddleware())

	if config.Cfg.HTTPAuth.Username != "" && config.Cfg.HTTPAuth.Password != "" {
		r.Use(BasicAuthMiddleware())
	}

	handlers.SetupRoutes(r)

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		filePath := strings.TrimPrefix(path, "/")
		if filePath == "" {
			filePath = "index.html"
		}

		file, err := subFS.Open(filePath)
		if err != nil {
			c.String(404, "Not Found: %s", path)
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			c.String(500, "Error getting file info")
			return
		}

		contentType := getMIMEType(filePath)
		c.Header("Content-Type", contentType)
		http.ServeContent(c.Writer, c.Request, "", stat.ModTime(), file.(io.ReadSeeker))
	})

	addr := fmt.Sprintf(":%d", listenPort)
	log.Printf("Starting Supervisord Monitor on %s", addr)
	log.Printf("Dashboard: http://localhost%s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, hasAuth := c.Request.BasicAuth()

		if !hasAuth || subtle.ConstantTimeCompare([]byte(user), []byte(config.Cfg.HTTPAuth.Username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(config.Cfg.HTTPAuth.Password)) != 1 {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
