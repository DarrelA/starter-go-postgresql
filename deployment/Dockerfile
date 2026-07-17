# Stage 1: Build the Go binary
FROM golang:1.26.5-trixie AS builder

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
# -o starter-go-postgresql: Specifies the output file name
RUN CGO_ENABLED=0 GOOS=linux go build -o starter-go-postgresql github.com/DarrelA/starter-go-postgresql/cmd/auth
RUN CGO_ENABLED=0 GOOS=linux go build -o starter-go-postgresql-migrate github.com/DarrelA/starter-go-postgresql/cmd/migrate

# Stage 2: Create a smaller image for the final binary
FROM alpine:3.24.1

# Add any necessary runtime dependencies for your application here
# RUN apk --no-cache add <dependencies>

# Copy the pre-built binary file from the builder stage
COPY --from=builder /docker_wd/starter-go-postgresql /root
COPY --from=builder /docker_wd/starter-go-postgresql-migrate /root

# Copy optional seed data; schema migrations are embedded in the binary.
COPY deployment/build/json /root/deployment/build/json

# Set the working directory for the container
WORKDIR /root

# Run the binary when the container starts
ENTRYPOINT ["./starter-go-postgresql"]
