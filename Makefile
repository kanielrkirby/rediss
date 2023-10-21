.PHONY: setup dev-setup prod-setup build-client build-server test-client test-server clean

# Build all

build: clean build-server build-client

# Clean all

clean: clean-server clean-client clean-lib


# Clean

clean-server:
	rm -rf ./server/tmp ./server/rediss-server ./server/rediss-server.exe ./server/rediss-server.log

clean-client:
	rm -rf ./client/tmp ./client/rediss-cli ./client/rediss-cli.exe ./client/rediss-cli.log

clean-lib:
	rm -rf ./lib/tmp ./lib/rediss-lib ./lib/rediss-lib.exe ./lib/rediss-lib.log


# Build

build-server:
	cd ./server && \
	go build -o ./rediss-server *.go

build-client:
	cd ./client && \
	go build -o ./rediss-cli *.go


# Dev

FLAGS :=

dev-server: setup-server-dev
	cd ./server && \
	go run *.go $(FLAGS)

dev-client: setup-client-dev
	cd ./client && \
	go run *.go $(FLAGS)
