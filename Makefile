
SHELL := /bin/bash

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

DB_DSN := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

ATLAS_ENV ?= local

PROTO_FILES = $(shell find app -name "*.proto" -not -path "*/google/*")
API_FILES   = $(shell find app -name "*.api")


.PHONY: help
help: 
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: debug
debug: 
	@echo "--- Debug Config ---"
	@echo "DB Host    : [$(DB_HOST)]"
	@echo "DB User    : [$(DB_USER)]"
	@echo "DB Name    : [$(DB_NAME)]"
	@echo "SSL Mode   : [$(DB_SSL_MODE)]"
	@echo "Full DSN   : $(DB_DSN)"
	@echo "Atlas Env  : [$(ATLAS_ENV)]"
	@echo "PayPal Mode: [$(PAYPAL_MODE)]"


.PHONY: diff
diff:
	@if [ -z "$(name)" ]; then echo "Error: Missing migration name. Please add name=migration_name"; exit 1; fi
	atlas migrate diff $(name) --env $(ATLAS_ENV)

.PHONY: apply
apply:
	atlas migrate apply --env $(ATLAS_ENV) --url "$(DB_DSN)"

.PHONY: down
down:
	@if [ -n "$(v)" ]; then \
		echo "Reverting to version: $(v) for service $(ATLAS_ENV)..."; \
		atlas migrate down --env $(ATLAS_ENV) --url "$(DB_DSN)" --to-version "$(v)"; \
	else \
		echo "Reverting the latest migration step for service $(ATLAS_ENV)..."; \
		atlas migrate down --env $(ATLAS_ENV) --url "$(DB_DSN)"; \
	fi

.PHONY: migrate-hash
migrate-hash:
	atlas migrate hash --env $(ATLAS_ENV)

.PHONY: status
status:
	atlas migrate status --env $(ATLAS_ENV) --url "$(DB_DSN)"

.PHONY: gen rpc gw mq

gen:
	@echo "1. Generating protobuf validation..."
	protoc -I . \
		--validate_out="lang=go,paths=source_relative:./dropshipbe" \
		dropshipbe.proto
		
	@echo "2. Generating gRPC source code (go-zero)......"
	goctl rpc protoc dropshipbe.proto --go_out=. --go-grpc_out=. --zrpc_out=.
	
	@echo "3. Generating Gateway descriptor..."
	protoc -I . --include_imports --descriptor_set_out=dropshipbe.pb dropshipbe.proto
	
	@echo "Done!"

rpc:
	@echo "Start gRPC server..."
	go run dropshipbe.go -f etc/dropshipbe.yaml

gw:
	@echo "Start API Gateway..."
	go run gateway/gateway.go -f etc/gateway.yaml

mq:
	@echo "Start Message Queue (Kafka Consumer)..."
	go run mq/mq.go -f etc/mq.yaml