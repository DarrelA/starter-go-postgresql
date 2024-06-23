# Define default values
APP_ENV ?= development

# Define default target
all: up

# Target to bring up the docker-compose services
up:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose up -d

# Target to bring down the docker-compose services
down:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose down

# Target to bring down the docker-compose services and named volumes
down-v:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose down -v

# Target to rebuild the docker-compose services
build:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose build

# Target to rebuild the docker-compose services
build-app:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose build app
