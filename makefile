# Define variables
APP_ENV ?= dev
DB ?= postgres redis
UI ?= pgAdmin

# Define default target
all: up

# Target to bring up the docker-compose services (excluding app-test)
up:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose up -d $(DB) $(UI) app

# Target to bring down the docker-compose services
d:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose down

# Target to bring down the docker-compose services and named volumes
dv:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose down -v

# Target to rebuild the docker-compose services
b:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose build

# Target to rebuild the docker-compose services (app and app-test)
bapp:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose build app app-test

# Target to rebuild the docker-compose services (app-test)
bat:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose build app-test

# Target to run tests (excluding app)
t:
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose up -d $(DB)
	@cd deployments && APP_ENV=$(APP_ENV) docker-compose run --rm app-test
	make dv

###############
# Maintenance #
###############

# Remove all dangling images and unused volumes
# If you want to skip the confirmation prompt, you can add the -f flag:
prune:
	@docker image prune -f
	@docker volume prune -f

# Update go.mod
umod:
	@go get -u
	@go mod tidy