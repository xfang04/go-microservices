FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerApp
MAIL_BINARY=mailerApp
LISTENER_BINARY=listenerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_logger build_mail build_listener
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker images..."
	docker-compose down
	@echo "Docker stopped!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Broker binary built!"

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ../authentication-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "auth binary built!"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ../logger-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "logger binary built!"

## build_mail: builds the mail binary as a linux executable
build_mail:
	@echo "Building mail binary..."
	cd ../mail-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
	@echo "mail binary built!"

## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo "Building listener binary..."
	cd ../listener-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "listener binary built!"

## front_end_linux: builds linux executable for front end
front_end:
	@echo "Building linux version of front end..."
	cd ../front-end && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o frontEnd ./cmd/web
	@echo "Done!"


## broker: stops broker-service, removes docker image, builds service, and starts it
broker: build_broker
	@echo "Building broker-service docker image..."
	- docker-compose stop broker-service
	- docker-compose rm -f broker-service
	docker-compose up --build -d broker-service
	docker-compose start broker-service
	@echo "broker-service rebuilt and started!"


## start: starts the front end
start:
	@echo "Starting front end"
	cd ../front-end && go build -o ${FRONT_END_BINARY} ./cmd/web
	cd ../front-end && ./${FRONT_END_BINARY}

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"

## test: runs all tests
test:
	@echo "Testing..."
	go test -v ./...

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@cd ../broker-service && rm -f ${BROKER_BINARY}
	@cd ../broker-service && go clean
	@cd ../front-end && go clean
	@cd ../front-end && rm -f ${FRONT_END_BINARY}
	@cd ../authentication-service && rm -f ${AUTH_BINARY}
	@cd ../authentication-service && go clean
	@echo "Cleaned!"
