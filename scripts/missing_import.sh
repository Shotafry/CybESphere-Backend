#!/bin/bash

echo "=== Go Module Info ==="
go mod tidy
go list -m

echo -e "\n=== Checking imports for dto package ==="
go list -f '{{.ImportPath}}: {{.Imports}}' ./internal/dto

echo -e "\n=== Checking imports for mappers package ==="
go list -f '{{.ImportPath}}: {{.Imports}}' ./internal/mappers

echo -e "\n=== Trying to build mappers package specifically ==="
go build ./internal/mappers

echo -e "\n=== Full build attempt with verbose output ==="
go build -v ./...

echo -e "\n=== Checking for circular dependencies ==="
go mod graph | grep -E "cybesphere-backend/internal/(dto|mappers)"