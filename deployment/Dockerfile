# Stage 1: Build the Go binary
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /docker_wd

# Copy the go.mod and go.sum files to download dependencies
# Copying these files separately allows Docker to cache the downloaded modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Remove the /testdata directory to keep the build context small
RUN rm -rf /docker_wd/testdata

# Build the application
# CGO_ENABLED=0: Disables CGO, which is appropriate if you want a pure Go binary without C dependencies
# GOOS=linux: Sets the target OS to Linux
# go build -a: Forces a rebuild of all packages and their dependencies,
#              regardless of whether they have changed or not
# -o starter-go-postgresql.test: Specifies the output file name
RUN CGO_ENABLED=0 GOOS=linux go build -a -o starter-go-postgresql github.com/DarrelA/starter-go-postgresql/cmd/auth

# Stage 2: Create a smaller image for the final binary
FROM alpine:latest

# Add any necessary runtime dependencies for your application here
# RUN apk --no-cache add <dependencies>

# Copy the pre-built binary file from the builder stage
COPY --from=builder /docker_wd/starter-go-postgresql /root

# Copy the necessary directories
COPY deployment/build/sql /root/deployment/build/sql
COPY deployment/build/json /root/deployment/build/json

# Set the working directory for the container
WORKDIR /root

# Run the binary when the container starts
ENTRYPOINT ["./starter-go-postgresql"]
