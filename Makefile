APP_NAME=Aqua.io
APP_ID=trading.aqua-io.app
VERSION=0.0.3
BUILD_NO=3
# if on Windows
ifeq ($(OS),Windows_NT)
BINARY_NAME="${APP_NAME}.exe"
MANIFEST_NAME=aqua-io.exe.manifest
else
# we're on a Mac
BINARY_NAME="${APP_NAME}.app"
endif

## build: build binary and package app
build:
ifeq ($(OS),Windows_NT)
	@del ${BINARY_NAME}
	fyne package -os windows -icon Icon.png -appID ${APP_ID} -appVersion ${VERSION} -appBuild ${BUILD_NO} -name ${APP_NAME} -release
else
	rm -rf ${BINARY_NAME}
	fyne package -os darwin -icon Icon.png -appID ${APP_ID} -appVersion ${VERSION} -appBuild ${BUILD_NO} -name ${APP_NAME} -release
endif

## cross compile the app for different architecture than the development machine (requires https://github.com/fyne-io/fyne-cross & docker)
cross-darwin-amd64:
	fyne-cross darwin -arch=amd64 -app-id=${APP_ID}
cross-windows-amd64:
	fyne-cross windows -arch=amd64 -app-id=${APP_ID}
	
## run: builds and runs the application
run:
ifeq ($(OS),Windows_NT)
	set DB_PATH=./sql.db && go run .\backend\internal\.
else
	env DB_PATH="./sql.db" go run .\backend\internal\.
endif

## debug: builds and runs the application in debug mode
debug:
ifeq ($(OS),Windows_NT)
	set DB_PATH=./sql.db && go run .\backend\internal\. -debug
else
	env DB_PATH="./sql.db" go run .\backend\internal\. -debug
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