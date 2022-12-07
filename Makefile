setup-env:
	@cp .env.example .env
	@docker-compose build --build-arg DEVELOPMENT=false iam

lint:
	@docker-compose run --rm iam go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run -v

test:
	@docker-compose run --rm iam go test ./...
