.PHONY: \
	build \
	static-build \
	run \
	tidy \
	lint \
	fmt \
	tests \
	version

build:
	go build -v

# Great for running in Docker containers.
static-build:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .

tools:
	go get golang.org/x/tools/cmd/goimports
	go get honnef.co/go/tools/cmd/staticcheck@latest
	go get github.com/kisielk/errcheck

run:
	go run .

tidy:
	go mod tidy -v

lint:
	go vet .
	errcheck .
	staticcheck .

fmt:
	go fmt .
	goimports -l -w .

tests:
	go clean -testcache .
	go test -cover -v .

version:
	@go version