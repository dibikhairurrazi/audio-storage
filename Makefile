.PHONY: mod.download
mod.download:
	GO111MODULE=on go mod download

.PHONY: tidy
tidy:
	GO111MODULE=on go mod tidy

.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor

.PHONY: generate-mock
generate-mock:
	bin/generate-mock.sh

build-http:
	@echo " >> building http binary"
	@go build -v -o audio-storage app/api/main.go

run-http: build-http
	@./audio-storage