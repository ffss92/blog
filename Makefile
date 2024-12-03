.PHONY: build
build:
	@go build -o bin/server cmd/server/*.go

.PHONY: run
run: build
	clear
	@bin/server

.PHONY: watch
watch:
	reflex -r '\.go$$' -d none -s make run
