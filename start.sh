#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

echo "========================================="
echo "  ADCMS - Starting All Services"
echo "========================================="
echo ""

# Start backend
echo "--- Starting Backend ---"
bash "$SCRIPT_DIR/api-backend/start.sh"
echo ""

# Start frontend
echo "--- Starting Frontend ---"
bash "$SCRIPT_DIR/api-frontend/start.sh"
echo ""

echo "========================================="
echo "  All services started!"
echo "  Backend API: http://localhost:8004"
echo "  Frontend:    http://localhost:3004"
echo "========================================="
