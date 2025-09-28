# https://just.systems

default:
    @just --list

# Check Go version and GOPATH
check-env:
    @echo "Checking Go version and GOPATH..."
    @go version
    @echo "GOPATH: ${GOPATH}"

# Run the application
run:
    @go run main.go

# build
build:
    @just fmt
    @echo "Building the application..."
    @go build -o bin/dem main.go

# Build for multiple platforms
build-all:
    @just fmt
    @echo "Building for macOS (amd64)..."
    @GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.GitCommit=$(git rev-parse HEAD)' -X 'main.GitBranch=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "main")'" -o bin/dem-mac-amd64 main.go
    @echo "Building for Windows (amd64)..."
    @GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.GitCommit=$(git rev-parse HEAD)' -X 'main.GitBranch=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "main")'" -o bin/dem-win-amd64.exe main.go
    @echo "Building for Linux (amd64)..."
    @GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.GitCommit=$(git rev-parse HEAD)' -X 'main.GitBranch=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "main")'" -o bin/dem-linux-amd64 main.go

# Format the entire project using goimports to organize imports
fmt:
    @echo "Formatting Go files and organizing imports..."
    @goimports -w .

# Pull the content from the remote repository to the local using the rebase
pull:
    @git fetch --all --tags --prune --jobs=10
    @git pull --rebase
    @just fmt

# Push the current branch to the branch of the remote repository.
push:
    @just fmt
    @git push origin main
    @git push gitee main
    @git push gitcode main

# Test: build first, then generate test data
test:
    @echo "Building the application..."
    @just build
    @echo "Generating test data..."
    @bash test/generate_test_data.sh

# Clean the build directory
clear:
    @echo "Clearing the build directory..."
    @rm -rf bin

# Clean the test data
clear-test:
    @echo "Clearing the test data..."
    @rm -rf ~/.dem/*.db

# Clear all
clear-all:
    @echo "Clearing all..."
    @just clear
    @just clear-test

# deploy mac
deploy-mac:
    @echo "Clearing all..."
    @just build
    @sudo cp ./bin/dem  /usr/local/bin/dem