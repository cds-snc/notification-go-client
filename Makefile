.PHONY: example-client test

example-client:
	@echo "Running example client..."
	@go run cmd/example_client/main.go


test:
	@go test -cover ./...