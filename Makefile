include .env
MIGRATIONS_PATH = ./cmd/migrate/migrations
VERSION := $(shell grep -m1 'const version' cmd/api/main.go | sed -E 's/.*"(.*)".*/\1/')

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: seed
seed: 
	@go run cmd/migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt

.PHONY: test
test:
	@go test ./... -v

.PHONY: docker-build
docker-build:
	@docker build -t social-api:$(VERSION) .

.PHONY: docker-run
docker-run: docker-build
	@docker run --rm -p 8080:8080 \
		-e DB_ADDR=postgres://admin:adminpassword@host.docker.internal/socialnetwork?sslmode=disable \
		-e REDIS_ENABLED=true \
		-e REDIS_ADDR=host.docker.internal:6379 \
		--name social-api \
		social-api:$(VERSION)

.PHONY: docker-full
docker-full:
	@docker compose up -d db redis
	@$(MAKE) docker-run

.PHONY: ecr-push
ecr-push: docker-build
	@ECR_REPO_URL=$$(cd terraform && terraform output -raw ecr_repository_url) && \
	ECR_REGISTRY=$${ECR_REPO_URL%%/*} && \
	AWS_REGION=$$(aws configure get region) && \
	aws ecr get-login-password --region $$AWS_REGION | docker login --username AWS --password-stdin $$ECR_REGISTRY && \
	docker tag social-api:$(VERSION) $$ECR_REPO_URL:$(VERSION) && \
	docker push $$ECR_REPO_URL:$(VERSION)