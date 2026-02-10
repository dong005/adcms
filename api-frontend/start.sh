#!/bin/bash

APP_NAME="adcms-frontend"
PID_FILE="./adcms-frontend.pid"
PORT=3004
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

export COREPACK_ENABLE_AUTO_PIN=0
export COREPACK_ENABLE_STRICT=0

echo "=== ADCMS Frontend ==="

# Check if already running
if [ -f "$PID_FILE" ]; then
    OLD_PID=$(cat "$PID_FILE")
    if kill -0 "$OLD_PID" 2>/dev/null; then
        echo "Stopping old process (PID: $OLD_PID)..."
        kill "$OLD_PID"
        sleep 2
        if kill -0 "$OLD_PID" 2>/dev/null; then
            kill -9 "$OLD_PID"
            sleep 1
        fi
        echo "Old process stopped."
    fi
    rm -f "$PID_FILE"
fi

# Also kill any process on the port
PID_ON_PORT=$(lsof -ti:$PORT 2>/dev/null)
if [ -n "$PID_ON_PORT" ]; then
    echo "Killing process on port $PORT (PID: $PID_ON_PORT)..."
    kill -9 $PID_ON_PORT 2>/dev/null
    sleep 1
fi

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "Installing dependencies..."
    pnpm install
fi

# Start dev server
echo "Starting frontend on port $PORT..."
nohup pnpm run dev:antd > /tmp/adcms-frontend.log 2>&1 &
NEW_PID=$!
echo "$NEW_PID" > "$PID_FILE"
echo "$APP_NAME started (PID: $NEW_PID)"
echo "Log: /tmp/adcms-frontend.log"
echo "URL: http://localhost:$PORT"
