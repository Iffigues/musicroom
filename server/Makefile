GOCMD=go

GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOINIT=$(GOCMD) mod init musicroom
GOTIDY=$(GOCMD) mod tidy
GOU= ./musicroom
all:  build
init:
	$(GOFMT)
	go mod init github.com/iffigues/musicroom
	go mod tidy
	$(GOGET)
build:
	$(GOFMT)
	go mod init github.com/iffigues/musicroom
	go mod tidy
	$(GOGET)
	$(GOBUILD)
fmt:
	$(GOFMT)
clean: 
	$(GOCLEAN)
	rm -rf $(GOU)
fclean:
	$(GOCLEAN)
	rm $(GOU) go.sum go.mod
	rm ./log/*
