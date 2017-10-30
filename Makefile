.PHONY: all vet mockgen test run build package clean install

APP_NAME=vpc-updater
APP_VERSION=0.0.1
APP_BUILD=`git log --pretty=format:'%h' -n 1`
GO_FLAGS= CGO_ENABLED=0
GO_LDFLAGS= -ldflags="-X main.AppVersion=$(APP_VERSION) -X main.AppName=$(APP_NAME) -X main.AppBuild=$(APP_BUILD)"
GO_BUILD_CMD=$(GO_FLAGS) go build $(GO_LDFLAGS)
BUILD_DIR=build
BINARY_NAME=$(APP_NAME)
MOCK_DIR=mocks

all: clean build package install

vet:
	@go vet

mockgen:
	@echo "Generating mocks..."
	mockgen -source=vendor/github.com/aws/aws-sdk-go/service/ec2/ec2iface/interface.go -destination=$(MOCK_DIR)/mock-ec2iface.go -package=mocks

test: mockgen
	@go test

install:
	cp $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 /usr/local/bin/$(BINARY_NAME)

run:
	go run ./*.go $(RUN_ARGS)

build: vet test
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64
	GOOS=darwin GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64
	GOOS=windows GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64

package:
	tar -C $(BUILD_DIR) -zcf $(BUILD_DIR)/$(BINARY_NAME)-$(APP_VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	tar -C $(BUILD_DIR) -zcf $(BUILD_DIR)/$(BINARY_NAME)-$(APP_VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	zip -q -j  $(BUILD_DIR)/$(BINARY_NAME)-$(APP_VERSION)-windows-amd64.zip $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64

clean:
	rm -Rf $(BUILD_DIR)s
