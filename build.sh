#!/bin/bash

echo "Building Supervisord Monitor..."

# Build frontend
echo "Building frontend..."
cd frontend
npm install
npm run build
cd ..

# Build backend
echo "Building backend..."
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

go build -ldflags "-s -w" -o supervisord-monitor

echo "Build complete!"
echo "Output: supervisord-monitor"
echo ""
echo "Usage:"
echo "  ./supervisord-monitor"
echo "  ./supervisord-monitor -config config.yaml"
echo "  ./supervisord-monitor -port 9000"