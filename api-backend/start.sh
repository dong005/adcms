#!/bin/bash

APP_NAME="adcms-api"
PID_FILE="./adcms-api.pid"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

echo "=== ADCMS API Backend ==="

# Check if already running
if [ -f "$PID_FILE" ]; then
    OLD_PID=$(cat "$PID_FILE")
    if kill -0 "$OLD_PID" 2>/dev/null; then
        echo "Stopping old process (PID: $OLD_PID)..."
        kill "$OLD_PID"
        sleep 2
        # Force kill if still running
        if kill -0 "$OLD_PID" 2>/dev/null; then
            kill -9 "$OLD_PID"
            sleep 1
        fi
        echo "Old process stopped."
    fi
    rm -f "$PID_FILE"
fi

# Build
echo "Building..."
go build -o "$APP_NAME" ./cmd/main.go
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi
echo "Build successful."

# Create logs directory
mkdir -p logs

# Start
echo "Starting $APP_NAME..."
nohup ./"$APP_NAME" > logs/stdout.log 2>&1 &
NEW_PID=$!
echo "$NEW_PID" > "$PID_FILE"
echo "$APP_NAME started (PID: $NEW_PID)"
echo "Log: $SCRIPT_DIR/logs/stdout.log"
