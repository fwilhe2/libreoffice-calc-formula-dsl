all: format build test

format:
	gofumpt -w $$(find . -name '*.go')

build:
	go build -v ./...

test:
	go test -v ./...

demo:
	./libreoffice-calc-formula-dsl