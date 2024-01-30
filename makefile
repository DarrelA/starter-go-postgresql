# Define default target
all: up

# Target to bring up the docker-compose services
up:
	@cd deployments && docker-compose up -d

# Target to bring down the docker-compose services
down:
	@cd deployments && docker-compose down

# Add more targets as needed
