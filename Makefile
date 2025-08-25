test:
	go test -v ./... -timeout 30s -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

.PHONY: test