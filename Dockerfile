# Stage 1: Build the Go binary
FROM golang:1.19.13-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# Stage 2: Create a minimal image with timezone support
FROM alpine:3.18.3 AS tzconfig

# Install tzdata to set the timezone
RUN apk add --no-cache tzdata

# Stage 3: Minimal image with the binary and timezone support
FROM scratch

# Set the timezone as a build argument (default is UTC)
ARG TZ="UTC"

# Copy the timezone data from the tzconfig stage
COPY --from=tzconfig /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=tzconfig /etc/passwd /etc/passwd

# Set the timezone in the container
ENV TZ=${TZ}

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/server /server

# Copy the configuration file
COPY --from=builder /app/config/config.yaml /config/config.yaml

# Expose the HTTP and gRPC ports
EXPOSE 8081 50001

# Set the working directory
WORKDIR /config

# Command to run the application
ENTRYPOINT ["/server"]
