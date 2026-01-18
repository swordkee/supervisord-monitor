package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"supervisord-monitor/config"
	"supervisord-monitor/services"

	"github.com/gin-gonic/gin"
)

type ProcessWithLog struct {
	services.ProcessInfo
	Log      string `json:"log"`
	HasError bool   `json:"has_error"`
}

func (p ProcessWithLog) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name        string `json:"name"`
		Group       string `json:"group"`
		State       int    `json:"state"`
		StateName   string `json:"statename"`
		Description string `json:"description"`
		Log         string `json:"log"`
		HasError    bool   `json:"has_error"`
	}{
		Name:        p.ProcessInfo.Name,
		Group:       p.ProcessInfo.Group,
		State:       p.ProcessInfo.State,
		StateName:   p.ProcessInfo.StateName,
		Description: p.ProcessInfo.Description,
		Log:         p.Log,
		HasError:    p.HasError,
	})
}

type ServerInfo struct {
	Name      string           `json:"name"`
	Version   string           `json:"version"`
	URL       string           `json:"url"`
	Processes []ProcessWithLog `json:"processes"`
	Error     string           `json:"error,omitempty"`
	HasAuth   bool             `json:"has_auth"`
}

func (s ServerInfo) MarshalJSON() ([]byte, error) {
	if s.Error != "" {
		return json.Marshal(struct {
			Name  string `json:"name"`
			Error string `json:"error"`
		}{
			Name:  s.Name,
			Error: s.Error,
		})
	}
	return json.Marshal(struct {
		Name      string           `json:"name"`
		Version   string           `json:"version"`
		URL       string           `json:"url"`
		Processes []ProcessWithLog `json:"processes"`
		HasAuth   bool             `json:"has_auth"`
	}{
		Name:      s.Name,
		Version:   s.Version,
		URL:       s.URL,
		Processes: s.Processes,
		HasAuth:   s.HasAuth,
	})
}

type DashboardResponse struct {
	Servers        []ServerInfo `json:"servers"`
	Refresh        int          `json:"refresh"`
	EnableAlarm    bool         `json:"enable_alarm"`
	SupervisorCols int          `json:"supervisor_cols"`
	ShowHost       bool         `json:"show_host"`
	Muted          bool         `json:"muted"`
}

func (r DashboardResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Servers        []ServerInfo `json:"servers"`
		Refresh        int          `json:"refresh"`
		EnableAlarm    bool         `json:"enable_alarm"`
		SupervisorCols int          `json:"supervisor_cols"`
		ShowHost       bool         `json:"show_host"`
		Muted          bool         `json:"muted"`
	}{
		Servers:        r.Servers,
		Refresh:        r.Refresh,
		EnableAlarm:    r.EnableAlarm,
		SupervisorCols: r.SupervisorCols,
		ShowHost:       r.ShowHost,
		Muted:          r.Muted,
	})
}

func GetDashboard(c *gin.Context) {
	muted := c.Query("mute") == "1"

	response := DashboardResponse{
		Refresh:        config.Cfg.Refresh,
		EnableAlarm:    config.Cfg.EnableAlarm,
		SupervisorCols: config.Cfg.SupervisorCols,
		ShowHost:       config.Cfg.ShowHost,
		Muted:          muted,
	}

	log.Printf("=== DEBUG: Starting GetDashboard ===")

	for name, serverCfg := range config.Cfg.SupervisorServers {
		log.Printf("Fetching data for server: %s", name)
		log.Printf("Server config: URL=%s, HasAuth=%v", serverCfg.URL, serverCfg.Username != "")

		client, err := services.NewSupervisorClient(name)
		if err != nil {
			log.Printf("ERROR: Failed to create client for %s: %v", name, err)
			response.Servers = append(response.Servers, ServerInfo{
				Name:  name,
				Error: err.Error(),
			})
			continue
		}

		version, err := client.GetSupervisorVersion()
		if err != nil {
			log.Printf("ERROR: Failed to get version for %s: %v", name, err)
		} else {
			log.Printf("Version for %s: %s", name, version)
		}

		processes, err := client.GetAllProcessInfo()
		if err != nil {
			log.Printf("ERROR: Failed to get processes for %s: %v", name, err)
			response.Servers = append(response.Servers, ServerInfo{
				Name:    name,
				Version: version,
				URL:     serverCfg.URL,
				Error:   err.Error(),
				HasAuth: serverCfg.Username != "",
			})
			continue
		}

		log.Printf("Got %d processes for server %s", len(processes), name)

		processWithLogs := make([]ProcessWithLog, 0, len(processes))
		for _, proc := range processes {
			log.Printf("Processing: Name=%s, Group=%s, State=%d", proc.Name, proc.Group, proc.State)
			logContent, _ := client.ReadProcessStderrLog(getProcessName(proc))
			processWithLogs = append(processWithLogs, ProcessWithLog{
				ProcessInfo: proc,
				Log:         logContent,
				HasError:    logContent != "",
			})
		}

		log.Printf("Built processWithLogs with %d items", len(processWithLogs))

		serverInfo := ServerInfo{
			Name:      name,
			Version:   version,
			URL:       serverCfg.URL,
			Processes: processWithLogs,
			HasAuth:   serverCfg.Username != "",
		}

		log.Printf("ServerInfo for %s: Name=%s, Processes len=%d, Version=%s",
			name, serverInfo.Name, len(serverInfo.Processes), serverInfo.Version)

		response.Servers = append(response.Servers, serverInfo)
	}

	log.Printf("=== DEBUG: Final Response ===")
	log.Printf("Total servers: %d", len(response.Servers))
	for i, srv := range response.Servers {
		log.Printf("Server %d: Name=%s, Processes len=%d, Error=%q, Version=%q",
			i, srv.Name, len(srv.Processes), srv.Error, srv.Version)
	}

	c.JSON(http.StatusOK, response)
}

func StartProcess(c *gin.Context) {
	server := c.Param("server")
	worker := c.Param("worker")

	client, err := services.NewSupervisorClient(server)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := client.StartProcess(worker); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func StopProcess(c *gin.Context) {
	server := c.Param("server")
	worker := c.Param("worker")

	client, err := services.NewSupervisorClient(server)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := client.StopProcess(worker); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func RestartProcess(c *gin.Context) {
	server := c.Param("server")
	worker := c.Param("worker")

	client, err := services.NewSupervisorClient(server)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := client.RestartProcess(worker); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func ClearProcessLog(c *gin.Context) {
	server := c.Param("server")
	worker := c.Param("worker")

	client, err := services.NewSupervisorClient(server)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := client.ClearProcessLogs(worker); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func StartAllProcesses(c *gin.Context) {
	server := c.Param("server")

	client, err := services.NewSupervisorClient(server)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := client.StartAllProcesses(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func StopAllProcesses(c *gin.Context) {
	server := c.Param("server")

	client, err := services.NewSupervisorClient(server)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := client.StopAllProcesses(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func RestartAllProcesses(c *gin.Context) {
	server := c.Param("server")

	client, err := services.NewSupervisorClient(server)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := client.RestartAllProcesses(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func getProcessName(proc services.ProcessInfo) string {
	if proc.Group != "" && proc.Group != proc.Name {
		return proc.Group + ":" + proc.Name
	}
	return proc.Name
}

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/dashboard", GetDashboard)
		api.POST("/start/:server/:worker", StartProcess)
		api.POST("/stop/:server/:worker", StopProcess)
		api.POST("/restart/:server/:worker", RestartProcess)
		api.POST("/clear/:server/:worker", ClearProcessLog)
		api.POST("/startall/:server", StartAllProcesses)
		api.POST("/stopall/:server", StopAllProcesses)
		api.POST("/restartall/:server", RestartAllProcesses)
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("%s %s completed in %v", c.Request.Method, c.Request.URL.Path, duration)
	}
}
