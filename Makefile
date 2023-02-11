setup-env:
	@cp .env.example .env
	@docker-compose build --build-arg DEVELOPMENT=false iam

lint:
	@docker-compose run --rm iam go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1 run -v

test:
	@docker-compose run --rm iam go test ./...

build-chart:
	cat VERSION | xargs -I {} helm package -u --version {} --app-version {} helm/iam
