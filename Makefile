TARGET ?= m5stick-c
TINYGO ?= tinygo

EXAMPLES := read led mode

.PHONY: agnostic test examples

test:
	go test ./...

examples:
	@for e in $(EXAMPLES); do \
		echo "=== examples/$$e ==="; \
		$(TINYGO) build -target=$(TARGET) -size short -o /tmp/m5unit-encoder-$$e.bin ./examples/$$e || exit 1; \
	done

# Board agnosticism + chip tag coverage, no hardware required.
agnostic:
	@rm -rf /tmp/enc-agn && mkdir -p /tmp/enc-agn && cd /tmp/enc-agn && \
	  MOD=$(shell awk '/^module /{print $$2}' go.mod) && \
	  printf 'module agn\n\ngo 1.26\n\nrequire %s v0.0.0\nreplace %s => %s\n' "$$MOD" "$$MOD" "$(CURDIR)" > go.mod && \
	  printf 'package main\n\nimport lib "%s"\n\nvar _ = lib.New\n\nfunc main() {}\n' "$$MOD" > main.go && \
	  go mod tidy >/dev/null && \
	  $(TINYGO) build -target=pico -o /tmp/enc-agn/pico.bin . && \
	  $(TINYGO) build -target=arduino-nano33 -o /tmp/enc-agn/nrf.bin . && \
	  printf 'package main\n\nimport lib "%s"\n\nvar _ = lib.NewI2C\n\nfunc main() {}\n' "$$MOD" > main.go && \
	  go mod tidy >/dev/null && \
	  for t in m5stick-c esp32-c3-devkit-rust-1 m5stamp-s3a pico; do \
	    $(TINYGO) build -target=$$t -o /tmp/enc-agn/$$t.bin . || exit 1; \
	  done && \
	  echo "agnostic: core pico+nrf OK, NewI2C esp32/c3/s3/rp2040 OK"
