#!/bin/bash
# Run Go linting with recommended configuration
# Usage: ./lint.sh [--fix]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Running Go linters...${NC}"

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo -e "${RED}golangci-lint is not installed.${NC}"
    echo "Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh"
    exit 1
fi

# Run go vet first (built-in)
echo -e "${YELLOW}Running go vet...${NC}"
go vet ./...

# Run gofmt check
echo -e "${YELLOW}Checking gofmt...${NC}"
if [ "$1" == "--fix" ]; then
    gofmt -w .
else
    FILES=$(gofmt -l .)
    if [ -n "$FILES" ]; then
        echo -e "${RED}Files not formatted:${NC}"
        echo "$FILES"
        echo -e "${YELLOW}Run with --fix or: gofmt -w .${NC}"
        exit 1
    fi
fi

# Run golangci-lint
echo -e "${YELLOW}Running golangci-lint...${NC}"
if [ "$1" == "--fix" ]; then
    golangci-lint run --fix ./...
else
    golangci-lint run ./...
fi

# Run go mod tidy check
echo -e "${YELLOW}Checking go.mod...${NC}"
go mod tidy
if [ -n "$(git diff --name-only go.mod go.sum 2>/dev/null)" ]; then
    echo -e "${RED}go.mod or go.sum changed. Please run 'go mod tidy' and commit.${NC}"
    exit 1
fi

echo -e "${GREEN}All checks passed!${NC}"
