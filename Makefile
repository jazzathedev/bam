run:
	go run ./cmd/bam

test:
	go test ./...

build:
	go build ./cmd/bam

run-shim:
	cd shim && go run ./

run-build:
	cd shim && go build ./