# if on Windows
ifeq ($(OS),Windows_NT)
BINARY_NAME=aqua-io.exe
MANIFEST_NAME=aqua-io.exe.manifest
else
# we're on a Mac
BINARY_NAME=aqua-io.app
endif
APP_NAME=Aqua.io
APP_ID=io.aqua-trading.app
VERSION=0.0.1
BUILD_NO=1

## build: build binary and package app
build:
ifeq ($(OS),Windows_NT)
	@del ${BINARY_NAME}
	fyne package -os windows -icon Icon.png -appID ${APP_ID} -appVersion ${VERSION} -appBuild ${BUILD_NO} -name ${APP_NAME} -release
else
	rm -rf ${BINARY_NAME}
	fyne package -os darwin -icon Icon.png -appID ${APP_ID} -appVersion ${VERSION} -appBuild ${BUILD_NO} -name ${APP_NAME} -release
endif
	
## run: builds and runs the application
run:
ifeq ($(OS),Windows_NT)
	set DB_PATH=./sql.db && go run .
else
	env DB_PATH="./sql.db" go run .
endif

## debug: builds and runs the application in debug mode
debug:
ifeq ($(OS),Windows_NT)
	set DB_PATH=./sql.db && go run . -debug
else
	env DB_PATH="./sql.db" go run . -debug
endif

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
ifeq ($(OS),Windows_NT)
	@del ${BINARY_NAME}
	@del ${MANIFEST_NAME}
else
	@rm -rf ${BINARY_NAME}
endif
	@echo "Cleaned!"

## test: runs all tests
test:
	go test -v ./...