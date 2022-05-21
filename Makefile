.PHONY: build run


build:
	go build -o survey-provider-service server.go

run:
	GIN_MODE=release ./survey-provider-service
