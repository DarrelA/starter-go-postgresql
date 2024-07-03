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
ENV_FILE=./configs/.env.$(APP_ENV)

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
	@cd deployments && $(VARS) docker-compose up -d $(DB) $(UI) app

# Target to bring down the docker-compose services
d:
	@cd deployments && $(VARS) docker-compose down

# Target to bring down the docker-compose services and named volumes
dv:
	@cd deployments && $(VARS) docker-compose down -v

# Target to rebuild the docker-compose app service
wa:
	@cd deployments && $(VARS) docker-compose build app

# Target to rebuild the docker-compose app-test service
wat:
	@cd deployments && $(VARS) docker-compose build app-test --build-arg APP_ENV=$(APP_ENV)

# Target to run tests (app-test)
t:
	@cd deployments && $(VARS) docker-compose run --rm app-test
	make dv

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