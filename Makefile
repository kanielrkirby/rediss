.PHONY: setup dev-setup prod-setup build-client build-server test-client test-server clean

# Build all

build: clean setup-build build-server build-client

# Clean all

clean: clean-server clean-client clean-lib

# Setup all build

setup-build: setup-server-build setup-client-build

# Setup all dev

setup-dev: setup-server-dev setup-client-dev


# Clean

clean-server:
	rm -rf ./server/tmp ./server/rediss-server ./server/rediss-server.exe ./server/rediss-server.log

clean-client:
	rm -rf ./client/tmp ./client/rediss-cli ./client/rediss-cli.exe ./client/rediss-cli.log

clean-lib:
	rm -rf ./lib/tmp ./lib/rediss-lib ./lib/rediss-lib.exe ./lib/rediss-lib.log


# Go mod

setup-server-dev:
	cd ./server && \
	go mod edit -replace github.com/piratey7007/rediss/lib=../lib && \
	go mod tidy

setup-client-dev:
	cd ./client && \
	go mod edit -replace github.com/piratey7007/rediss/lib=../lib && \
	go mod tidy

setup-server-build:
	cd ./server && \
	go mod edit -dropreplace github.com/piratey7007/rediss/lib && \
	go mod tidy

setup-client-build:
	cd ./client && \
	go mod edit -dropreplace github.com/piratey7007/rediss/lib && \
	go mod tidy


# Build

build-server:
	cd ./server && \
	go build -o ./rediss-server *.go

build-client:
	cd ./client && \
	go build -o ./rediss-cli *.go


# Dev

dev-server: setup-server-dev
	cd ./server && \
	go run *.go

dev-client: setup-client-dev
	cd ./client && \
	go run *.go
