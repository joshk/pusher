EXAMPLES := $(wildcard examples/*.go)

examples: */**.go
	for example in $(EXAMPLES); do \
		go run $$example; \
	done

fmt: */**.go
	go fmt ./...

test: */**.go
	go test ./...
