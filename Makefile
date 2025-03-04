BUILD_DIR 	:= build
SCRIPTS_DIR := scripts
MSG ?= "Default commit message"

build:
	go build -o $(BUILD_DIR)/API main/*

run:
	cd $(BUILD_DIR) && ./API

runserver: build run

format:
	./$(SCRIPTS_DIR)/pre-commit

commit: format
	git add .
	git commit -m "$(MSG)"

push: commit
	git push

clean:
	rm -rf $(BUILD_DIR)/*

.PHONY: build run clean