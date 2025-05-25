PROJECT_ID := $(shell cd infra/terraform-cloudrun-api && terraform output -raw project_id 2>/dev/null)
REGION := $(shell cd infra/terraform-cloudrun-api && terraform output -raw region 2>/dev/null)
REPOSITORY_URL := $(shell cd infra/terraform-cloudrun-api && terraform output -raw repository_url 2>/dev/null)
IMAGE_NAME := cms-api
IMAGE_TAG := latest

# Terraform操作
plan:
	cd infra/terraform-cloudrun-api && terraform plan

apply:
	cd infra/terraform-cloudrun-api && terraform apply

destroy:
	cd infra/terraform-cloudrun-api && terraform destroy

# Docker操作
docker-login:
	gcloud auth configure-docker $(REGION)-docker.pkg.dev

docker-build:
	cd src && docker build --target prd -t $(REPOSITORY_URL)/$(IMAGE_NAME):$(IMAGE_TAG) .

docker-push: docker-login docker-build
	docker push $(REPOSITORY_URL)/$(IMAGE_NAME):$(IMAGE_TAG)
