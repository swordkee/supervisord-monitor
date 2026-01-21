# Supervisord Monitor (Go + Vue)

A multi-server Supervisord monitoring tool built with Gin + Vue, with all static files embedded into a single executable.

## Features

- 🎯 Real-time monitoring of multiple Supervisord servers
- 🚀 Start/Stop/Restart individual or all processes
- 📊 Real-time process status and uptime display
- 🔔 Error log viewing and sound alerts
- 🔇 Mute function (using browser local storage)
- 🔄 Auto-refresh page (configurable)
- 🔐 Support for authenticated and non-authenticated Supervisord servers
- 🔒 Web interface HTTP Basic Authentication (optional)
- 📦 Single executable file, no additional deployment required
- ⚡ High-performance Go backend + Vue 3 frontend

## Tech Stack

### Backend
- **Gin**: High-performance Go Web framework
- **XML-RPC**: Communication with Supervisord RPC2 interface
- **Embed**: Embed static files into binary
- **Viper**: Configuration management

### Frontend
- **Vue 3**: Progressive JavaScript framework
- **Vite**: Next-generation frontend build tool
- **Axios**: HTTP client
- **Bootstrap 5**: Responsive UI framework

## Installation

### Prerequisites

- Go 1.21+
- Node.js 16+ (required for build only)
- Supervisord server (RPC2 interface must be enabled)

### Configure Supervisord Server

Enable `inet_http_server` in each Supervisord configuration file:

```ini
[inet_http_server]
port=*:9001
username="yourusername"
password="yourpass"
```

Restart Supervisord service:

```bash
sudo supervisorctl restart all
```

### Build Project

#### Windows

```bash
build.bat
```

#### Linux/Mac

```bash
chmod +x build.sh
./build.sh
```

After build, a single executable will be generated:
- Windows: `supervisord-monitor.exe`
- Linux/Mac: `supervisord-monitor`

## Configuration

Edit `config.yaml` to configure monitored servers:

```yaml
supervisor_cols: 2      # Dashboard columns (2 or 3)
refresh: 10             # Refresh interval (seconds), 0 to disable
enable_alarm: true      # Enable/disable alarm sound
show_host: false        # Show hostname after server name
timeout: 3              # RPC2 interface connection timeout (seconds)
port: 8080              # Web service port

# Web interface HTTP Basic Authentication (optional)
# Leave empty to disable authentication
http_auth:
  username: "admin"     # Login username
  password: "admin123"  # Login password

supervisor_servers:
  - name: "server01"
    url: "http://server01.app/RPC2"
    port: "9001"
    username: "yourusername"
    password: "yourpass"
  - name: "server02"
    url: "http://server02.app/RPC2"
    port: "9001"
```

## Usage

### Basic Run

```bash
# Windows
supervisord-monitor.exe

# Linux/Mac
./supervisord-monitor
```

### Specify Config File

```bash
supervisord-monitor.exe -config /path/to/config.yaml
```

### Specify Port

```bash
supervisord-monitor.exe -port 9000
```

### Combined Usage

```bash
supervisord-monitor.exe -config config.yaml -port 8080
```

Default access: `http://localhost:8080`

### HTTP Authentication

If `http_auth` is configured in `config.yaml`, a login dialog will appear:

- Username: `http_auth.username` from config
- Password: `http_auth.password` from config

To disable authentication, leave `username` and `password` empty.

## Features

### Process Control

- **Start**: Start individual process
- **Stop**: Stop individual process
- **Restart**: Restart individual process (stop then start)
- **Start All**: Start all processes on server
- **Stop All**: Stop all processes on server
- **Restart All**: Restart all processes on server

### Process Status

- **RUNNING**: Process is running (green)
- **STOPPED**: Process is stopped (black)
- **STARTING**: Process is starting (blue)
- **FATAL**: Process failed to start (red, triggers alarm)
- **EXITED**: Process has exited (red, triggers alarm)

### Error Log

- Warning icon appears next to process name when stderr log exists
- Click warning icon to view error log content
- Alarm triggers automatically when status is FATAL

### Sound Alert

- Alert sound plays when FATAL status or stderr log is detected
- Click "Mute" button to mute (saved to browser local storage)
- Mute state persists after page refresh

### Auto Refresh

- Auto-refresh page every 10 seconds by default
- Adjust interval by modifying `refresh` option in config
- Set to 0 to disable auto-refresh

## API Endpoints

### Get Dashboard Data

```
GET /api/dashboard?mute=1
```

### Process Control

```
POST /api/start/{server}/{worker}
POST /api/stop/{server}/{worker}
POST /api/restart/{server}/{worker}
POST /api/clear/{server}/{worker}
POST /api/startall/{server}
POST /api/stopall/{server}
POST /api/restartall/{server}
```

## Troubleshooting

### "Did not receive a '200 OK' response from remote server"

- Check firewall and network connection
- Ensure Supervisord RPC2 interface is accessible
- Verify URL and port configuration are correct

### "HTTP/1.0 401 Unauthorized"

- Incorrect username or password
- Check authentication info in config file

### "UNKNOWN_METHOD"

- Supervisord v3+ requires RPC interface to be enabled
- Add `[rpcinterface:supervisor]` section in config file

### Cannot Start Executable

- Ensure executable has execute permission (Linux/Mac)
- Windows may need to allow unsigned applications

## Development

### Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install Node dependencies
cd frontend
npm install
```

### Frontend Development

```bash
cd frontend
npm run dev
```

Access `http://localhost:5173`

### Backend Development

```bash
go run main.go
```

Access `http://localhost:8080`

## Project Structure

```
supervisord-monitor-go/
├── config/
│   └── config.go          # Configuration management
├── handlers/
│   └── handlers.go        # API handlers
├── services/
│   └── supervisord.go     # Supervisord XML-RPC client
├── embed.go               # Static file embedding
├── main.go                # Main entry point
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   └── ServerCard.vue  # Server card component
│   │   ├── App.vue        # Main app component
│   │   ├── main.js        # Entry point
│   │   └── style.css      # Global styles
│   ├── index.html         # HTML template
│   ├── package.json       # Node dependencies
│   └── vite.config.js     # Vite configuration
├── sounds/
│   └── alert.mp3          # Alert sound
├── config.yaml            # Configuration file
├── build.bat              # Windows build script
└── build.sh               # Linux/Mac build script
```

## License

MIT License

## Credits

- Original project: [mlazarov/supervisord-monitor](https://github.com/mlazarov/supervisord-monitor)