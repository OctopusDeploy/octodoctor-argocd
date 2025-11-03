BIN_NAME=argocd-doctor

# Get version from git tag or hash, default to "dev"
VERSION := $(shell git describe --tags --exact-match HEAD 2>/dev/null || git rev-parse --short HEAD 2>/dev/null || echo "dev")

argocd-doctor-build: 
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -v -ldflags "-X main.version=$(VERSION)" -o ${BIN_NAME} .

.PHONY: lint
lint:
	golangci-lint run --fix

.PHONY: test
test:
	go test -v ./...

.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

.PHONY: coverage-html
coverage-html: test-coverage
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: coverage
coverage:
	go test -cover ./...


.PHONY: release
release:
	make BIN_NAME=argocd-doctor-darwin-amd64 GOOS=darwin GOARCH=amd64 argocd-doctor-build
	make BIN_NAME=argocd-doctor-darwin-arm64 GOOS=darwin GOARCH=arm64 argocd-doctor-build
	make BIN_NAME=argocd-doctor-linux-amd64 GOOS=linux GOARCH=amd64 argocd-doctor-build
	make BIN_NAME=argocd-doctor-linux-arm64 GOOS=linux GOARCH=arm64 argocd-doctor-build
	make BIN_NAME=argocd-doctor-windows-amd64.exe GOOS=windows argocd-doctor-build

.PHONY: clean
clean:
	rm -f $(BIN_NAME)
	rm -rf coverage.out coverage.html
	rm -rf argocd-doctor-output