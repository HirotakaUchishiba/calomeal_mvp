#!/bin/bash

# Protocol Buffers code generation script
# This script generates Go code from .proto files

set -e

echo "🔧 Generating Protocol Buffers code..."

# Create output directories
mkdir -p proto/foods/v1
mkdir -p proto/logs/v1
mkdir -p proto/analytics/v1

# Generate Go code for foods service
echo "📦 Generating foods service..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/foods/v1/foods.proto

# Generate Go code for logs service
echo "📦 Generating logs service..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/logs/v1/logs.proto

# Generate Go code for analytics service
echo "📦 Generating analytics service..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/analytics/v1/analytics.proto

echo "✅ Protocol Buffers code generation completed!"
echo ""
echo "Generated files:"
echo "- proto/foods/v1/foods.pb.go"
echo "- proto/foods/v1/foods_grpc.pb.go"
echo "- proto/logs/v1/logs.pb.go"
echo "- proto/logs/v1/logs_grpc.pb.go"
echo "- proto/analytics/v1/analytics.pb.go"
echo "- proto/analytics/v1/analytics_grpc.pb.go"
