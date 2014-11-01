all:
	scripts/build.sh

dist:
	scripts/dist.sh

qa: test vet fmt lint

lint:
	go get github.com/golang/lint/golint
	$(GOPATH)/bin/golint

fmt:
	gofmt -s -w . model controller

clean:
	rm -f bin/* || true
	rm -rf .gopath || true

test:
	go test -v ./...

vet:
	go vet ./...

here: build qa

build:
	go build -o bin/app-scheduler

.PHONY: all	dist clean test
