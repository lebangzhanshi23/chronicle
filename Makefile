.PHONY: build run clean test frontend-build

BINARY_NAME=chronicle

# Git info
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null)
GIT_DATE=$(shell git log -1 --format=%cs 2>/dev/null)
BUILD_TIME=$(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null)

LDFLAGS=-X main.gitCommit=$(GIT_COMMIT) -X main.gitDate=$(GIT_DATE) -X main.buildTime=$(BUILD_TIME) -X github.com/yuyudeqiu/chronicle/internal/handler.GitCommit=$(GIT_COMMIT) -X github.com/yuyudeqiu/chronicle/internal/handler.GitDate=$(GIT_DATE) -X github.com/yuyudeqiu/chronicle/internal/handler.BuildTime=$(BUILD_TIME)

build: frontend-build
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) main.go
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME).exe main.go

run:
	go run main.go server

test:
	go test -v ./...

clean:
	rm -rf bin/
	rm -rf frontend/dist/

frontend-build:
	cd frontend && npm install && npm run build

help:
	@echo "Usage:"
	@echo "  make build           - Build frontend and backend"
	@echo "  make run             - Run the server"
	@echo "  make test            - Run tests"
	@echo "  make clean           - Clean builds"
	@echo "  make frontend-build  - Build frontend only"
