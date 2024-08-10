# Dockerfile

# Use a base image with Go installed
FROM quay.io/projectquay/golang:1.20 as builder

# Define build arguments for target OS and architecture
ARG targetos=linux
ARG targetarch=amd64

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the entire context into the working directory
COPY . .

# Run the make build command with the specified target OS and architecture
RUN make build TARGETOS=$targetos TARGETARCH=$targetarch

# Use a minimal base image for the final stage
FROM alpine:latest

# Set the working directory for the final image
WORKDIR /

# Copy the built binary from the builder stage
COPY --from=builder /go/src/app/kbot .
# Copy the SSL certificates
# COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs

RUN apk add --no-cache ca-certificates

# Set the entrypoint to the built binary
CMD ["./kbot", "start"]

