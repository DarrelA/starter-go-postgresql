# Define variables
APP_ENV ?= dev
DEPENDENCIES ?= postgres pgAdmin redis

# Define default target
all: up

# Target to bring up the docker-compose services (excluding app-test)
up:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose up -d $(DEPENDENCIES) app

# Target to bring down the docker-compose services
down:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose down

# Target to bring down the docker-compose services and named volumes
down-v:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose down -v

# Target to rebuild the docker-compose services
build:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose build

# Target to rebuild the docker-compose services (app and app-test)
build-app:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose build app app-test

# Target to run tests (excluding app)
t:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose up -d $(DEPENDENCIES)
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose run --rm app-test
	make down