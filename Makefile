OUTPUT = bin
RELEASER = goreleaser
APP_NAME 	 ?= terrafactor


.PHONY: test
test:
	go test ./...

.PHONY: go-build
go-build:
	go build -o ${OUTPUT}/$(APP_NAME) main.go

.PHONY: clean
clean:
	@test ! -e ${OUTPUT}/${BIN_NAME} || rm ${OUTPUT}/${BIN_NAME}

.PHONY: install
install:
	go get ./...

main: main.go
	goreleaser build --snapshot --rm-dist
