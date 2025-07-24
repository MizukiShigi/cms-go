PROJECT_ID := $(shell cd infra/terraform-cloudrun-api && terraform output -raw project_id 2>/dev/null)
REGION := $(shell cd infra/terraform-cloudrun-api && terraform output -raw region 2>/dev/null)
REPOSITORY_URL := $(shell cd infra/terraform-cloudrun-api && terraform output -raw repository_url 2>/dev/null)
IMAGE_NAME := cms-api
IMAGE_TAG := latest

# Google Cloud操作
login:
	gcloud auth application-default login

# Terraform操作
plan:
	cd infra/terraform-cloudrun-api && terraform plan

apply:
	cd infra/terraform-cloudrun-api && terraform apply

destroy:
	cd infra/terraform-cloudrun-api && terraform destroy

# ローカル専用
apply-image:
	cd infra/terraform-cloudrun-api && terraform apply -var="container_image=$(REPOSITORY_URL)/$(IMAGE_NAME):$(IMAGE_TAG)"  -replace=google_cloud_run_service.main

# Docker操作
docker-login:
	gcloud auth configure-docker $(REGION)-docker.pkg.dev

docker-build:
	cd src && docker build --platform linux/amd64 --target prd -t $(REPOSITORY_URL)/$(IMAGE_NAME):$(IMAGE_TAG) .

docker-push: docker-login docker-build
	docker push $(REPOSITORY_URL)/$(IMAGE_NAME):$(IMAGE_TAG)

# go tool
test:
	cd src && go test $$(go list ./... | grep -v '/infrastructure/db/sqlboiler')

fmt:
	cd src && go fmt ./...

vet:
	cd src && go vet ./...

# DBマイグレーション
migration:
	cd src && sqlboiler psql

# モック生成
# 使用例: make mock-repo REPO=post_repository
mock-repo:
	@if [ -z "$(REPO)" ]; then \
		echo "使用例: make mock-repo REPO=post_repository"; \
		exit 1; \
	fi
	cd src && go run go.uber.org/mock/mockgen@latest -source=internal/domain/repository/$(REPO).go -destination=mocks/repository/mock_$(REPO).go -package=repository

mock-all:
	cd src && go run go.uber.org/mock/mockgen@latest -source=internal/domain/repository/post_repository.go -destination=mocks/repository/mock_post_repository.go -package=repository
	cd src && go run go.uber.org/mock/mockgen@latest -source=internal/domain/repository/user_repository.go -destination=mocks/repository/mock_user_repository.go -package=repository
	cd src && go run go.uber.org/mock/mockgen@latest -source=internal/domain/repository/tag_repository.go -destination=mocks/repository/mock_tag_repository.go -package=repository
	cd src && go run go.uber.org/mock/mockgen@latest -source=internal/domain/repository/transaction_manager.go -destination=mocks/repository/mock_transaction_manager.go -package=repository
	cd src && go run go.uber.org/mock/mockgen@latest -source=internal/domain/repository/image_repository.go -destination=mocks/repository/mock_image_repository.go -package=repository

# 下位互換のため
mock: mock-all
