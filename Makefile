BUILD_DIR 	:= build
SCRIPTS_DIR := scripts
MSG ?= "Default commit message"

build:
	go build -o $(BUILD_DIR)/API main/*

run: build generate-docs
	cd $(BUILD_DIR) && ./API

runserver: build run

runRedisServer :
	$(RUN_REDIS_SERVER_COMMAND)

runDBServer :
	$(RUN_DB_SERVER_COMMAND)

format:
	./$(SCRIPTS_DIR)/pre-commit

commit: format
	echo "\033[32mstagging changes...\033[0m"
	@git add .
	echo "\033[32mcommiting changes...\033[0m"
	@git commit -m "$(MSG)"

push: commit
	echo "\033[32mpushing to remote repo...\033[0m"
	@git push

generate-docs:
	swag init -g main/main.go -o ./docs 

clean:
	rm -rf $(BUILD_DIR)/*

.PHONY: build run clean
