#!/bin/bash

# Development script for the game server

set -e

echo "🎮 Starting Game Server Development Environment"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Start dependencies (MySQL, Redis)
echo "📦 Starting dependencies..."
docker-compose up -d mysql redis

# Wait for dependencies to be ready
echo "⏳ Waiting for dependencies to be ready..."
sleep 10

# Check if dependencies are healthy
echo "🔍 Checking dependencies health..."
docker-compose ps mysql redis

# Install dependencies
echo "📥 Installing Go dependencies..."
go mod download
go mod tidy

# Generate protocol files if needed
echo "🔧 Generating protocol files..."
if command -v flatc &> /dev/null; then
    make flatbuffers
else
    echo "⚠️  FlatBuffers compiler not found. Install it to generate .fbs files"
fi

# Start the server with hot reload
echo "🚀 Starting server with hot reload..."
if command -v air &> /dev/null; then
    air
else
    echo "⚠️  Air not found. Installing..."
    go install github.com/cosmtrek/air@latest
    air
fi
