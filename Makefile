.PHONY: gen dev mocks test

gen:
	go generate ./...

dev:
	go run main.go

mocks:
	mockery --dir=store --name=Store

test: gen
	go test ./...
