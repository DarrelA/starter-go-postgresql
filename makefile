# Define default target
all: up

# Target to bring up the docker-compose services
up:
	@cd deployments && docker-compose up -d

# Target to bring down the docker-compose services
down:
	@cd deployments && docker-compose down

# Target to bring down the docker-compose services and named volumes
down-v:
	@cd deployments && docker-compose down -v

# Target to rebuild the docker-compose services
build:
	@cd deployments && docker-compose build

# Target to rebuild the docker-compose services
build-app:
	@cd deployments && docker-compose build app

