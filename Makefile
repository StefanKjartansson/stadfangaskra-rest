all: test build

NAME = stadfangaskra

ifeq ($(OS),Windows_NT)
	OUTNAME := $(NAME)-win
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		OUTNAME := $(NAME)-linux
	endif
	ifeq ($(UNAME_S),Darwin)
		OUTNAME := $(NAME)-osx
	endif
endif

deps:
	go get github.com/djimenez/iconv-go
	go get github.com/StefanKjartansson/isnet93
	go get code.google.com/p/gorilla/mux

test: .PHONY
	go test -v 

build:
	go build -o ${OUTNAME} bin/stadfangaskra.go	

.PHONY:
