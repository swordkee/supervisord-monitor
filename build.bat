@echo off
echo Building Supervisord Monitor...

REM Build frontend
echo Building frontend...
cd frontend
call npm install
call npm run build
cd ..

REM Build backend
echo Building backend...
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1

go build -ldflags "-s -w" -o supervisord-monitor.exe

echo Build complete!
echo Output: supervisord-monitor.exe
echo.
echo Usage:
echo   supervisord-monitor.exe
echo   supervisord-monitor.exe -config config.yaml
echo   supervisord-monitor.exe -port 9000
