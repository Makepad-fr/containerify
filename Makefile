BUILD_TARGET ?= unknown
EXECUTABLE_NAME ?= unknown
EXECUTABLE_PATH ?= unknown
VALID_TARGETS = cli server
VALID_ENVS = docker local
MAIN_FILE_NAME = main.go
BUILD_ENV = local
IMAGE_BASE_NAME = containerify
export CGO_ENABLED=0
export GOOS=linux

.PHONY: build
build: install-dependencies
ifeq ($(filter $(BUILD_TARGET),$(VALID_TARGETS)),)
	$(error Unsupported BUILD_TARGET: $(BUILD_TARGET))
endif
ifndef EXECUTABLE_NAME
	$(error EXECUTABLE_NAME is not set)
endif
ifndef EXECUTABLE_PATH
	$(error EXECUTABLE_PATH is not set)
endif
ifeq ($(filter ${BUILD_ENV}, ${VALID_ENVS}),)
	$(error Unsupported BUILD_ENV ${BUILD_ENV})
endif
ifeq ($(BUILD_ENV),local)
	@$(MAKE) create-executable-path
	@go build -o ${EXECUTABLE_PATH}/${EXECUTABLE_NAME} ./${BUILD_TARGET}/${MAIN_FILE_NAME}
endif
ifeq ($(BUILD_ENV), docker)
	docker build --build-arg BUILD_TARGET=$(BUILD_TARGET) -t ${IMAGE_BASE_NAME}-${BUILD_TARGET} .
endif
.PHONY: install-dependencies
install-dependencies:
	@echo "Sync workspace"
	@go work sync
.PHONY: create-executable-path
create-executable-path:
	@mkdir -p ${EXECUTABLE_PATH}