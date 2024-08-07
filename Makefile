# Makefile

# Get the application name from the git remote URL
APP=$(shell basename $(shell git remote get-url origin))
# Docker registry
REGISTRY=us-west1-docker.pkg.dev
PROJECT_ID=kbot-429800
REPO_NAME=devops-repo
# Get the version from the latest git tag and short commit hash
VERSION=$(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
# Default target OS and architecture
TARGETOS=linux
TARGETARCH=amd64

# Format the Go code
format:
	gofmt -s -w ./

# Display the latest version
latestVersion:
	echo ${VERSION}

# Lint the Go code
lint:
	golint

# Run tests
test:
	go test -v

# Get the dependencies
get:
	go get

# Build the Go binary with the specified OS and architecture
build: get latestVersion format
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o kbot -ldflags "-X="github.com/andrey-ilin/kbot/cmd.appVersion=${VERSION}

# Build the Docker image with the specified build arguments
image:
	@echo "*****************************************************"
	@echo "Build for ${TARGETOS} with ${TARGETARCH} architecture"
	@echo "*****************************************************"
	docker build --build-arg targetos=${TARGETOS} --build-arg targetarch=${TARGETARCH} . -t ${REGISTRY}/${PROJECT_ID}/${REPO_NAME}/${APP}:${VERSION}-${TARGETOS}-${TARGETARCH}


# Push the Docker image to the registry
push:
	docker push ${REGISTRY}/${PROJECT_ID}/${REPO_NAME}/${APP}:${VERSION}-${TARGETOS}-${TARGETARCH}
	
# Clean up the built binary
clean:
	rm -rf kbot
	@echo "*****************************************************"
	@echo "Removing ${REGISTRY}/${PROJECT_ID}/${REPO_NAME}/${APP}:${VERSION}-${TARGETOS}-${TARGETARCH} image"
	@echo "*****************************************************"
	@images=$$(docker images --filter=reference="${REGISTRY}/${PROJECT_ID}/${REPO_NAME}/${APP}:${VERSION}-*" -q); \
	if [ -z "$$images" ]; then \
		echo "No images found for ${REGISTRY}/${APP}:${VERSION}-*"; \
	else \
		docker rmi $$images -f; \
	fi

# Build for Linux
linux:
	@$(MAKE) build TARGETOS=linux

# Build for macOS
macos:
	@$(MAKE) build TARGETOS=darwin

# Build for Windows
windows:
	@$(MAKE) build TARGETOS=windows

