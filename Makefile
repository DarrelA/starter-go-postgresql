#####################
#  Define Variables #
#####################

APP_ENV ?= dev
DB = postgres redis
UI = pgAdmin
VARS = APP_ENV=$(APP_ENV) POSTGRES_DB=$(POSTGRES_DB) POSTGRES_USER=$(POSTGRES_USER) 

#####################
#    Env Configs    #
#####################

# Path to the environment-specific .env file
ENV_FILE=./internal/infrastructure/config/.env.$(APP_ENV)

# Check if the environment-specific .env file exists
ifeq (,$(wildcard $(ENV_FILE)))
  $(error "$(ENV_FILE) file not found")
endif

# Include environment-specific variables from .env.${APP_ENV}
include $(ENV_FILE)
export $(shell sed 's/=.*//' $(ENV_FILE))

#####################
#    make <cmd>     #
#####################

# Define default target
all: up

# Target to bring up the docker-compose services (excluding app-test)
up:
	@cd deployment && $(VARS) docker-compose up -d $(DB) $(UI) app

# Target to bring down the docker-compose services
d:
	@cd deployment && APP_ENV=$(APP_ENV) docker-compose down
	make lg

# Target to bring down the docker-compose services, named volumes, and remove unused containers
dv:
	@cd deployment && APP_ENV=$(APP_ENV) docker-compose down -v
	@docker container prune -f
	make lg

# Target to rebuild the docker-compose app service
wa:
	@cd deployment && $(VARS) docker-compose build app

# Target to rebuild the docker-compose app-unit-test service & run test
ut:
	@cd deployment && $(VARS) docker-compose build app-unit-test --build-arg APP_ENV=$(APP_ENV)
	@cd deployment && $(VARS) docker-compose run app-unit-test
	make dv

# Target to rebuild the docker-compose app-integration-test service & run test
it:
	@cd deployment && $(VARS) docker-compose build app-integration-test --build-arg APP_ENV=$(APP_ENV)
	@cd deployment && $(VARS) docker-compose run app-integration-test
	make dv


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

# Update go.mod
umod:
	@go get -u
	@go mod tidy