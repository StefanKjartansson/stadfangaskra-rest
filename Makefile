all: test

test: .PHONY
	go test -v ./...

.PHONY:
