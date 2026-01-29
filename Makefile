.PHONY: run vendor

run: vendor
	go run main.go

vendor:
	go mod tidy
	go mod vendor
