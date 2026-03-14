include .env

MIGRATE_POSTGRES_DIR ?= ./migrations/postgres
MIGRATE_CLICKHOUSE_DIR ?= ./migrations/clickhouse

DB_CONFIGS := pg|$(MIGRATE_POSTGRES_DIR)|postgres|$(POSTGRES_DB_DSN) \
			  ch|$(MIGRATE_CLICKHOUSE_DIR)|clickhouse|$(CLICKHOUSE_DB_DSN)

help::
	@echo "Usage: make [target]"
	@echo ""

define make-migrate-targets
$(eval alias := $(word 1,$(subst |, ,$(1))))
$(eval dir := $(word 2,$(subst |, ,$(1))))
$(eval driver := $(word 3,$(subst |, ,$(1))))
$(eval dsn := $(word 4,$(subst |, ,$(1))))

migrate-$(alias)-up: ; goose -dir $(dir) $(driver) $(dsn) up
migrate-$(alias)-down: ; goose -dir $(dir) $(driver) "$(dsn)" down
migrate-$(alias)-status: ; goose -dir $(dir) $(driver) "$(dsn)" status
migrate-$(alias)-create: ; goose -dir $(dir) $(driver) "$(dsn)" create $(name) sql

.PHONY: migrate-$(alias)-up migrate-$(alias)-down migrate-$(alias)-status migrate-$(alias)-create

help::
	@echo "  migrate-$(alias)-up       Run all pending $(driver) migrations"
	@echo "  migrate-$(alias)-down     Roll back last $(driver) migration"
	@echo "  migrate-$(alias)-status   Show $(driver) migration status"
	@echo "  migrate-$(alias)-create [name=migration_name] Create a new $(driver) migration"
	@echo ""
endef

$(foreach config,$(DB_CONFIGS),$(eval $(call make-migrate-targets,$(config))))

.DEFAULT_GOAL := help
