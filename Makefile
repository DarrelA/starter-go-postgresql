#####################
#  Define Variables #
#####################

# Test targets use the test configuration unless APP_ENV is explicitly set.
ifneq (,$(filter ut it,$(MAKECMDGOALS)))
APP_ENV ?= test
else
APP_ENV ?= dev
endif
DB := postgres redis
UI := pgadmin
COMPOSE := docker compose
VARS := APP_ENV=$(APP_ENV)
DOCS_ADDR ?= localhost:6060

.DEFAULT_GOAL := all
.PHONY: all init-env rotate-keys up up-pgadmin migrate d dv wa ut it lg e prune umod docs

#####################
#    Env Configs    #
#####################

# Path to the environment-specific .env file
ENV_FILE := ./internal/infrastructure/config/.env.$(APP_ENV)
ENV_EXAMPLE := $(ENV_FILE).example
PGADMIN_SERVER_FILE := ./internal/infrastructure/db/postgres/json/servers.$(APP_ENV).json
PGADMIN_SERVER_EXAMPLE := ./internal/infrastructure/db/postgres/json/servers.$(APP_ENV).example.json
ENV_REQUIRED_GOALS := all rotate-keys up up-pgadmin migrate d dv wa ut it lg e
REQUESTED_GOALS := $(if $(MAKECMDGOALS),$(MAKECMDGOALS),all)

# Only runtime targets require local environment configuration. Maintenance and
# bootstrap targets remain usable in a fresh clone.
ifneq (,$(filter $(ENV_REQUIRED_GOALS),$(REQUESTED_GOALS)))
ifeq (,$(wildcard $(ENV_FILE)))
  $(error "$(ENV_FILE) file not found; run 'make init-env APP_ENV=$(APP_ENV)'")
endif
include $(ENV_FILE)
export $(shell sed -n 's/^\([A-Za-z_][A-Za-z0-9_]*\)=.*/\1/p' $(ENV_FILE))
else
-include $(ENV_FILE)
endif

#####################
#    make <cmd>     #
#####################

# Define default target
all: up

# Create local configuration without committing credentials.
init-env:
	@test -f "$(ENV_EXAMPLE)" || { echo "$(ENV_EXAMPLE) file not found"; exit 1; }
	@test -f "$(PGADMIN_SERVER_EXAMPLE)" || { echo "$(PGADMIN_SERVER_EXAMPLE) file not found"; exit 1; }
	@if [ ! -f "$(ENV_FILE)" ]; then \
		cp "$(ENV_EXAMPLE)" "$(ENV_FILE)"; \
		deployment/build/scripts/refresh_token_keygen.sh "$(ENV_FILE)"; \
		echo "Created $(ENV_FILE) with generated JWT keys"; \
	else \
		echo "Keeping existing $(ENV_FILE)"; \
	fi
	@if [ ! -f "$(PGADMIN_SERVER_FILE)" ]; then \
		cp "$(PGADMIN_SERVER_EXAMPLE)" "$(PGADMIN_SERVER_FILE)"; \
		echo "Created $(PGADMIN_SERVER_FILE)"; \
	else \
		echo "Keeping existing $(PGADMIN_SERVER_FILE)"; \
	fi

# Generate a new local access/refresh signing-key pair.
rotate-keys:
	@deployment/build/scripts/refresh_token_keygen.sh "$(ENV_FILE)"
	@echo "Rotated JWT keys in $(ENV_FILE)"

# Start the application and its required data stores without pgAdmin.
up:
	@cd deployment && $(VARS) $(COMPOSE) up -d $(DB) app

# Start the application and its data stores with the optional pgAdmin UI.
up-pgadmin:
	@cd deployment && $(VARS) $(COMPOSE) --profile admin up -d $(DB) $(UI) app

# Run the one-shot migration job and return its exit status.
migrate:
	@cd deployment && $(VARS) $(COMPOSE) --profile release build migrate
	@cd deployment && $(VARS) $(COMPOSE) --profile release run --rm migrate

# Target to bring down the docker-compose services
d:
	@cd deployment && APP_ENV=$(APP_ENV) $(COMPOSE) --profile admin down
	@$(MAKE) lg

# Target to bring down the Compose services and their named volumes.
dv:
	@cd deployment && APP_ENV=$(APP_ENV) $(COMPOSE) --profile admin down -v
	@$(MAKE) lg

# Target to rebuild the docker-compose app service
wa:
	@cd deployment && $(VARS) $(COMPOSE) build app

# Rebuild the unit-test image, run it, and clean up even when tests fail.
ut:
	@status=0; \
	cd deployment && \
	$(VARS) $(COMPOSE) build app-unit-test --build-arg APP_ENV=$(APP_ENV) && \
	$(VARS) $(COMPOSE) run --rm app-unit-test || status=$$?; \
	cd ..; \
	$(MAKE) APP_ENV=$(APP_ENV) dv || cleanup_status=$$?; \
	if [ "$$status" -ne 0 ]; then exit "$$status"; fi; \
	exit "$${cleanup_status:-0}"

# Rebuild the integration-test image, run it, and clean up even when tests fail.
it:
	@status=0; \
	cd deployment && \
	$(VARS) $(COMPOSE) build app-integration-test --build-arg APP_ENV=$(APP_ENV) && \
	$(VARS) $(COMPOSE) run --rm app-integration-test || status=$$?; \
	cd ..; \
	if [ "$$status" -eq 0 ]; then \
		go tool cover -html=./testdata/reports/covdatafiles/coverage.out \
			-o ./testdata/reports/it_coverage.html || status=$$?; \
	fi; \
	$(MAKE) APP_ENV=$(APP_ENV) dv || cleanup_status=$$?; \
	if [ "$$status" -ne 0 ]; then exit "$$status"; fi; \
	exit "$${cleanup_status:-0}"


# Format log
lg:
	@deployment/build/scripts/format_app_log.sh

# Echo variables
e:
	@echo "$(VARS)"


#####################
#    Maintenance    #
#####################

# Remove all dangling images and unused volumes
# If you want to skip the confirmation prompt, you can add the -f flag:
prune:
	@docker image prune -f
	@docker volume prune -f

# Update dependencies used by all packages, then remove unused requirements.
umod:
	@go get -u ./...
	@go mod tidy

# Serve this module's Go documentation without requiring a global pkgsite installation.
docs:
	@go tool pkgsite -open -http=$(DOCS_ADDR) -list=false .
