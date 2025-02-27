BUILD_DIR := build

build:
	go build -o $(BUILD_DIR)/API main/main.go main/routes.go

run:
	cd $(BUILD_DIR) && ./API

clean:
	rm -rf $(BUILD_DIR)/*

.PHONY: build run clean