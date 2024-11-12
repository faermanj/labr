# docker build -f Containerfile --no-cache --progress=plain -t $REGISTRY_USERNAME/labr:$(cat version.txt) .
# docker run -it --rm $REGISTRY_USERNAME/labr:$(cat version.txt)
# docker login -u $REGISTRY_USERNAME -p $REGISTRY_PASSWORD
# docker push $REGISTRY_USERNAME/labr:$(cat version.txt)

# Use a Go base image for building the application
FROM golang:alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download all dependencies. Dependencies are cached if the go.mod and go.sum are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN GOOS=linux GOARCH=amd64 go build -o labr .

# Start a new stage from the scratch image (smaller final image)
FROM alpine:latest  

# Install bash (required for running .sh scripts inside the container)
RUN apk --no-cache add bash

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/labr .

# Command to run the binary
ENTRYPOINT ["./labr"]
