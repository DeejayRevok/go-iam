setup-env:
	@cp .env.example .env
	@docker-compose build --build-arg DEVELOPMENT=false uaa

lint:
	@docker-compose run --rm uaa go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run -v

test:
	@docker-compose run --rm uaa go test ./...
