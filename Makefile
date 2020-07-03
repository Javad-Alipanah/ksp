SERVER=bin/server
STATIC_SERVER=bin/server_static
SOURCES=$(shell find . -name '*.go' -not -name '*_test.go')

all: $(SERVER)

static: $(STATIC_SERVER)

format:
	find . -name '*.go' -not -path "./.cache/*" | xargs -n1 go fmt

check: format
	git diff
	git diff-index --quiet HEAD

lint:
	golangci-lint run --skip-dirs=test --deadline 3m0s

test:
	go test -cover ./... -coverprofile .coverage.txt
	cat .coverage.txt | grep "/internal\|mode:" > .coverage.pkg
	go tool cover -func .coverage.pkg

clean:
	rm -rf bin

bin/%: ./%.go $(SOURCES)
	go build -o $@ $<
#	strip -s $@

bin/%_static: ./%.go $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $@ $<
#	strip -s $@

bin:
	mkdir -p $@

$(TEST_DIR):
	mkdir -p $@

.PHONY: all static format check lint test clean
