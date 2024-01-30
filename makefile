# Define default target
all: up

# Target to bring up the docker-compose services
up:
	@cd deployments && docker-compose up -d

# Target to bring down the docker-compose services
down:
	@cd deployments && docker-compose down

# Target to rebuild the docker-compose services
rebuild:
	@cd deployments && docker-compose build
