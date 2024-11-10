.phony: build
build:
	go build -o broadcast ./cmd/.

start:
	./broadcast start

connect:
	./broadcast connect