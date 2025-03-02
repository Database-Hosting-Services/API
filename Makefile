BUILD_DIR := build

build:
	go build -o $(BUILD_DIR)/API main/*

run:
	cd $(BUILD_DIR) && ./API

runserver: build run

clean:
	rm -rf $(BUILD_DIR)/*

.PHONY: build run clean