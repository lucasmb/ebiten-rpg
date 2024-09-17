
OUTPUT='./gorpg'

PHONY: build
build: 
	go build -o $(OUTPUT)

PHONY: run
run:
	go run *.go