.PHONY: test test-ci coverage build generate sbom

test:
	go test ./...

test-ci:
	mkdir -p build
	go tool gotestsum --format testname --junitfile build/test-results.xml --jsonfile build/test-results.json -- -coverprofile=build/coverage.out ./...
	go tool cover -html=build/coverage.out -o build/coverage.html

build:
	go build ./...

generate:
	go generate ./...

sbom:
	mkdir -p build
	go tool cyclonedx-gomod mod -json -output build/bom.json .
