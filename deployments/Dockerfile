# Start from the latest Golang base image for building the binary
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies
# Copying these files separately allows Docker to cache the downloaded modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
# -o flag sets the output binary name
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o starter-go-postgresql .

# Start from a smaller base image for the final image
FROM alpine:latest  

# Add any necessary runtime dependencies for your application here
# RUN apk --no-cache add <dependencies>

# Copy the pre-built binary file from the builder stage
COPY --from=builder /app/starter-go-postgresql /root/

# Set the working directory for the container
WORKDIR /root/

# Run the binary when the container starts
CMD ["./starter-go-postgresql"]
