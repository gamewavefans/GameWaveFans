
OUTPUT_DIR=bin
COVER_FILE=cover.out

.PHONY: build test clean all

all: build_native test vet cover

build_all: build_linux_32 build_linux_64 build_windows_32 build_windows_64

build_native:
	mkdir -p ${OUTPUT_DIR}/native
	go build  -ldflags "-s -w" -o ${OUTPUT_DIR}/native ./cmd/...

build_linux_32:
	mkdir -p ${OUTPUT_DIR}/linux_32
	GOOS=linux GOARCH=386 go build  -ldflags "-s -w" -o ${OUTPUT_DIR}/linux_32 ./cmd/...

build_linux_64:
	mkdir -p ${OUTPUT_DIR}/linux_64
	GOOS=linux GOARCH=amd64 go build  -ldflags "-s -w" -o ${OUTPUT_DIR}/linux_64 ./cmd/...

build_windows_32:
	mkdir -p ${OUTPUT_DIR}/windows_32
	GOOS=windows GOARCH=386 go build  -ldflags "-s -w" -o ${OUTPUT_DIR}/windows_32 ./cmd/...

build_windows_64:
	mkdir -p ${OUTPUT_DIR}/windows_64
	GOOS=windows GOARCH=amd64 go build  -ldflags "-s -w" -o ${OUTPUT_DIR}/windows_64 ./cmd/...

test:
	go test ./cmd/... ./pkg/...

clean:
	rm -rf ${COVER_FILE}
	rm -rf ${OUTPUT_DIR}


cover:
	go test -coverprofile ${COVER_FILE} ./...

cover_show: cover
	go tool cover -html=${COVER_FILE}

vet:
	go vet ./cmd/... ./pkg/...

fmt:
	go fmt ./cmd/... ./pkg/...	
