# service names
ServiceLauncher := Launcher
ServiceISW := ISW
ServiceAPI := API

BUILD_DIR ?= build
CMD_DIR := cmd

clean:
	@rm -rf $(BUILD_DIR)

build: build-isw build-launcher build-api copy-weather copy-config

copy-config:
	@mkdir -p $(BUILD_DIR)
	@cp ./config.yaml $(BUILD_DIR)

build-launcher:
	@mkdir -p $(BUILD_DIR)
	@go build -o "$(ServiceLauncher)" ./$(CMD_DIR)/launcher

build-api:
	@mkdir -p $(BUILD_DIR)
	@go build -o "$(ServiceAPI)" ./$(CMD_DIR)/api

build-isw:
	@mkdir -p $(BUILD_DIR)
	@go build -o "$(ServiceISW)" ./$(CMD_DIR)/isw

copy-weather:
	@mkdir -p $(BUILD_DIR)
	@cp -r tools/weather $(BUILD_DIR)
