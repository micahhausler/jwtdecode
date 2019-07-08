COVERAGE=cover.out
BINARY=jwtdecode

$(COVERAGE):
	go test -v -coverprofile=$(COVERAGE) ./pkg
test: $(COVERAGE)

cover: test
	go tool cover -func=$(COVERAGE)

$(BINARY): test
	go build -o $(BINARY) main.go

build: $(BINARY)

clean:
	rm -rf ./$(COVERAGE) ./$(BINARY)

.PHONY: test cover build clean
