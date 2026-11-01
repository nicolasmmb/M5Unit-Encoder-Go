TARGET ?= m5stick-c
TINYGO ?= tinygo

EXAMPLES := read led mode

.PHONY: test examples

test:
	go test ./...

examples:
	@for e in $(EXAMPLES); do \
		echo "=== examples/$$e ==="; \
		$(TINYGO) build -target=$(TARGET) -size short -o /tmp/m5unit-encoder-$$e.bin ./examples/$$e || exit 1; \
	done
