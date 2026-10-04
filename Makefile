.PHONY: check test vet fmt

check: fmt vet test

fmt:
	@echo "==> gofmt"
	@files=$$(find . -name '*.go' -not -path './.git/*'); \
	unformatted=$$(gofmt -l $$files); \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt needed on:"; echo "$$unformatted"; exit 1; \
	fi

vet:
	@echo "==> go vet"
	go vet ./...

test:
	@echo "==> go test"
	go test ./...
