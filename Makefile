# Makefile for building the golang package
REGISTRY ?= ghcr.io
REPONAME ?= toversus
IMAGE_BASE ?= otel-demo
BACKEND_DOCKER_TAG ?= backend
FRONTEND_DOCKER_TAG ?= frontend

.PHONY: docker
docker: docker-backend docker-frontend

.PHONY: docker-backend
docker-backend:
	docker buildx build -t $(REGISTRY)/$(REPONAME)/$(IMAGE_BASE):$(BACKEND_DOCKER_TAG) .

.PHONY: docker-frontend
docker-frontend:
	docker buildx build -t $(REGISTRY)/$(REPONAME)/$(IMAGE_BASE):$(FRONTEND_DOCKER_TAG) -f ./frontend/Dockerfile .
