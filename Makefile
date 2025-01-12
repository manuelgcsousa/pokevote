BIN_DIR=bin

.PHONY: default seed clean

default:
	@echo "Available commands:"
	@echo "  seed  ~> Run the seed program"
	@echo "  build ~> Build the pokevote binary"
	@echo "  run   ~> Run the pokevote program"
	@echo "  clean ~> Clean build artifacts"

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

seed:
	go run cmd/seed/main.go

build: $(BIN_DIR)
	go build -o $(BIN_DIR)/pokevote cmd/pokevote/main.go

run: build
	go run cmd/pokevote/main.go

clean:
	go clean
	rm -rf $(BIN_DIR)
