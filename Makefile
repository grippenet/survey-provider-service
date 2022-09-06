.PHONY: build run build-release

build:
	go build -o survey-provider-service server.go

run:
	GIN_MODE=release ./survey-provider-service

build-release:
	goreleaser release --skip-publish --rm-dist