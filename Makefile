BUILD_DIR=./build
TARGET=$(BUILD_DIR)/marvin

all: clean test build resources

test:
	go test -v ./...

clean:
	@if [ -d $(BUILD_DIR) ]; then rm -rf $(BUILD_DIR); fi;
	go clean

pre-build:
	@mkdir -p $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)/config
	go mod vendor

.PHONY: build
build:
	make pre-build
	if [ -f $(TARGET) ]; then rm $(TARGET); fi;
	go build -v -o $(TARGET)

build-rpi:
	make pre-build
	if [ -f $(TARGET)_rpi ]; then rm $(TARGET)_rpi; fi;
	GOOS=linux GOARCH=arm GOARM=6 go build -v -o $(TARGET)_rpi

.PHONY: resources
resources:
	make pre-build
	@test -f config.json && cp -v config.json $(BUILD_DIR)/config
	@test -f endpoints.json && cp -v endpoints.json $(BUILD_DIR)/config
	@cp -rv webapp $(BUILD_DIR)
	@cp -rv resources $(BUILD_DIR)

deploy-rpi:
ifeq ($(RPI_HOST),)
	@echo "use option RPI_HOST=<rpi-hostname> to deploy to rpi"
	@exit 1
endif
	make clean
	make build-rpi
	make resources
	@mv $(TARGET)_rpi $(TARGET)
	@echo "stopping service on $(RPI_HOST)..."
	ssh pi@${RPI_HOST} sudo systemctl stop marvin
	sleep 10
	scp -r $(BUILD_DIR)/* pi@${RPI_HOST}:/home/pi/marvin/
	@echo "starting service on $(RPI_HOST)..."
	ssh pi@${RPI_HOST} sudo systemctl start marvin
	make clean

