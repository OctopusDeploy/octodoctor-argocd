BIN_NAME=argocd-doctor

argocd-doctor-build: 
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -v -o ${BIN_NAME} .

.PHONY: lint
lint:
	golangci-lint run --fix

.PHONY: test
test:
	go test -v 


.PHONY: release
release:
	make BIN_NAME=argocd-doctor-darwin-amd64 GOOS=darwin GOARCH=amd64 argocd-doctor-build
	make BIN_NAME=argocd-doctor-darwin-arm64 GOOS=darwin GOARCH=arm64 argocd-doctor-build
	make BIN_NAME=argocd-doctor-linux-amd64 GOOS=linux GOARCH=amd64 argocd-doctor-build
	make BIN_NAME=argocd-doctor-linux-arm64 GOOS=linux GOARCH=arm64 argocd-doctor-build
	make BIN_NAME=argocd-doctor-windows-amd64.exe GOOS=windows argocd-doctor-build
