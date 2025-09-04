APP_NAME = minio-server
BIN_DIR = bin
BIN_PATH = $(BIN_DIR)/$(APP_NAME)
CONFIG_PATH = /home/gautam/Desktop/DevHub/GoHub/Lemma/minio-server/resource/config/config.yml

# Ensure bin directory exists
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# Build the binary
build: $(BIN_DIR)
	go build -o $(BIN_PATH) main.go

# Run the binary (after building it)
run: build
	$(BIN_PATH) --config $(CONFIG_PATH)

# Clean the build
clean:
	rm -rf $(BIN_DIR)
