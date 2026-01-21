package main

import (
	"crypto/subtle"
	"flag"
	"fmt"
	"log"
	"net/http"

	"supervisord-monitor/config"
	"supervisord-monitor/handlers"

	"github.com/gin-gonic/gin"
)

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

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(handlers.LoggerMiddleware())

	if config.Cfg.HTTPAuth.Username != "" && config.Cfg.HTTPAuth.Password != "" {
		r.Use(BasicAuthMiddleware())
	}

	handlers.SetupRoutes(r)

	addr := fmt.Sprintf(":%d", listenPort)
	log.Printf("Starting Supervisord Monitor API on %s", addr)

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