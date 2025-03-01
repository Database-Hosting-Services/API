BUILD_DIR := build

build:
	go build -o $(BUILD_DIR)/API main/*

run:
	cd $(BUILD_DIR) && ./API

runserver:
	go build -o $(BUILD_DIR)/API main/*
	cd $(BUILD_DIR) && ./API

clean:
	rm -rf $(BUILD_DIR)/*

.PHONY: build run clean