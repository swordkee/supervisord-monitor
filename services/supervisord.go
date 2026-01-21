package services

import (
	"fmt"
	"log"
	"time"

	"supervisord-monitor/config"

	supervisor "github.com/abrander/go-supervisord"
)

type SupervisorClient struct {
	server string
	config *config.SupervisorServer
	client *supervisor.Client
}

type ProcessInfo struct {
	Name        string `json:"name"`
	Group       string `json:"group"`
	State       int    `json:"state"`
	StateName   string `json:"statename"`
	Description string `json:"description"`
}

func NewSupervisorClient(serverName string) (*SupervisorClient, error) {
	cfg, err := config.GetServerConfig(serverName)
	if err != nil {
		return nil, err
	}

	log.Printf("Creating Supervisord client for %s at %s", serverName, cfg.URL)

	opts := []supervisor.ClientOption{}

	if cfg.Username != "" && cfg.Password != "" {
		opts = append(opts, supervisor.WithAuthentication(cfg.Username, cfg.Password))
	}

	client, err := supervisor.NewClient(cfg.URL, opts...)
	if err != nil {
		return nil, err
	}

	return &SupervisorClient{
		server: serverName,
		config: cfg,
		client: client,
	}, nil
}

func formatProcessDescription(proc *supervisor.ProcessInfo) string {
	log.Printf("formatProcessDescription: State=%d, StateName=%s, Start=%d, Now=%d, Pid=%d",
		proc.State, proc.StateName, proc.Start, proc.Now, proc.Pid)

	if proc.State == supervisor.StateRunning {
		if proc.Pid > 0 && proc.Start > 0 {
			now := proc.Now
			if now == 0 {
				now = int(time.Now().Unix())
			}

			duration := time.Duration(now-proc.Start) * time.Second
			hours := int(duration.Hours())
			minutes := int(duration.Minutes()) % 60

			var uptime string
			if hours > 24 {
				days := hours / 24
				uptime = fmt.Sprintf("uptime %dd %dh %dm", days, hours%24, minutes)
			} else if hours > 0 {
				uptime = fmt.Sprintf("uptime %dh %dm", hours, minutes)
			} else {
				uptime = fmt.Sprintf("uptime %dm", minutes)
			}

			return fmt.Sprintf("pid %d, %s", proc.Pid, uptime)
		}

		if proc.Pid > 0 {
			return fmt.Sprintf("pid %d", proc.Pid)
		}
	}

	return proc.StateName
}

func (s *SupervisorClient) GetAllProcessInfo() ([]ProcessInfo, error) {
	processes, err := s.client.GetAllProcessInfo()
	if err != nil {
		log.Printf("Failed to get processes: %v", err)
		return nil, err
	}

	result := make([]ProcessInfo, len(processes))
	for i, proc := range processes {
		result[i] = ProcessInfo{
			Name:        proc.Name,
			Group:       proc.Group,
			State:       int(proc.State),
			StateName:   proc.StateName,
			Description: formatProcessDescription(&proc),
		}
		log.Printf("Process %d: %+v", i, result[i])
	}

	log.Printf("Successfully parsed %d processes", len(result))
	return result, nil
}

func (s *SupervisorClient) GetSupervisorVersion() (string, error) {
	return s.client.GetSupervisorVersion()
}

func (s *SupervisorClient) StartProcess(name string) error {
	return s.client.StartProcess(name, true)
}

func (s *SupervisorClient) StopProcess(name string) error {
	return s.client.StopProcess(name, true)
}

func (s *SupervisorClient) RestartProcess(name string) error {
	return s.client.StartProcess(name, true)
}

func (s *SupervisorClient) ClearProcessLogs(name string) error {
	return s.client.ClearProcessLogs(name)
}

func (s *SupervisorClient) ReadProcessStderrLog(name string) (string, error) {
	return s.client.ReadProcessStderrLog(name, -50000, 0)
}

func (s *SupervisorClient) StartAllProcesses() error {
	_, err := s.client.StartAllProcesses(true)
	return err
}

func (s *SupervisorClient) StopAllProcesses() error {
	_, err := s.client.StopAllProcesses(true)
	return err
}

func (s *SupervisorClient) RestartAllProcesses() error {
	_, err := s.client.StopAllProcesses(true)
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	_, err = s.client.StartAllProcesses(true)
	return err
}

func getProcessName(proc ProcessInfo) string {
	if proc.Group != "" && proc.Group != proc.Name {
		return proc.Group + ":" + proc.Name
	}
	return proc.Name
}
