GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=hpg
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_UNIX=$(BINARY_NAME)_unix
EXECUTABLE = $(BINARY_WINDOWS)

all: test build
build: 
		$(GOBUILD) -o $(EXECUTABLE) -v
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_WINDOWS)
		rm -f $(BINARY_UNIX)
run:
		$(GOBUILD) -o $(EXECUTABLE) -v ./...
		./$(EXECUTABLE)
deps:
		$(GOGET) github.com/markbates/goth
		$(GOGET) github.com/markbates/pop


# Cross compilation
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v