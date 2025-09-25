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

# Build for multiple platforms
build:
    @echo "Building for macOS (amd64)..."
    @GOOS=darwin GOARCH=amd64 go build -o bin/dem-mac-amd64 main.go
    @echo "Building for Windows (amd64)..."
    @GOOS=windows GOARCH=amd64 go build -o bin/dem-win-amd64.exe main.go
    @echo "Building for Linux (amd64)..."
    @GOOS=linux GOARCH=amd64 go build -o bin/dem-linux-amd64 main.go

# Format the entire project using goimports to organize imports
fmt:
    @echo "Formatting Go files and organizing imports..."
    @goimports -w .

# Pull the content from the remote repository to the local using the rebase
pull:
    @git fetch --all --tags --prune --jobs=10
    @git pull --rebase

# Push the current branch to the branch of the remote repository.
push:
    @git push origin main
    @git push gitee main
    @git push gitcode main
