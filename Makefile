lint:
	golangci-lint run --go=1.18

test:
	go test ./...