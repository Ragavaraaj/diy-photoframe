.DEFAULT_GOAL := run

.PHONY: fmt vet build brun

fmt:
	echo "Formatting code..."
	go fmt ./...

vet:
	echo "Running go vet..."
	go vet ./...

build:
	echo "Building the application..."
	go build -o diy-photoframe ./build

brun: build
	echo "Running the application in background..."
	./build/diy-photoframe 
	
run: 
	echo "Running the application..."
	go run main.go