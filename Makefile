# if on Windows
ifeq ($(OS),Windows_NT)
BINARY_NAME=Copy IO.exe
else
# we're on a Mac
BINARY_NAME=Copy IO.app
endif
APP_NAME=Copy.io
VERSION=1.0.1
BUILD_NO=1

## build: build binary and package copy_io
build:
ifeq ($(OS),Windows_NT)
	@del ${BINARY_NAME}
else
	rm -rf ${BINARY_NAME}
endif
	fyne package -appVersion ${VERSION} -appBuild ${BUILD_NO} -name ${APP_NAME} -release

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
else
	@rm -rf ${BINARY_NAME}
endif
	@echo "Cleaned!"

## test: runs all tests
test:
	go test -v ./...