APP_NAME=app
BUILD_DIR=bin
MAIN=./cmd
PUBLISHER_TOOL_NAME=publisher
PUBLISHER_TOOL=./tools/publisher/mqtt_publisher.go

CLIENT_TOOL_NAME=client
CLIENT_TOOL=./tools/client/ws_client.go

.PHONY: run build clean

include .env
export $(shell sed 's/=.*//' .env)

build:
	go mod tidy
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN)

run: build
	$(BUILD_DIR)/$(APP_NAME)

clean:
	rm -rf $(BUILD_DIR)

build-publisher:
	go mod tidy
	go build -o $(BUILD_DIR)/$(PUBLISHER_TOOL_NAME) $(PUBLISHER_TOOL)
	
run-publisher: build-publisher
	$(BUILD_DIR)/$(PUBLISHER_TOOL_NAME)

build-client:
	go mod tidy
	go build -o $(BUILD_DIR)/$(CLIENT_TOOL_NAME) $(CLIENT_TOOL)
	
build-client:
	go mod tidy
	go build -o $(BUILD_DIR)/$(CLIENT_TOOL_NAME) $(CLIENT_TOOL)

run-client: build-client
	$(BUILD_DIR)/$(CLIENT_TOOL_NAME) $(ARGS)
