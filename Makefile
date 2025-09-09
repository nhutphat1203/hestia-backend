APP_NAME=app
BUILD_DIR=bin
MAIN=./cmd

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
