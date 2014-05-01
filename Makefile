GO ?= go

all: build

build:
		$(GO) build

clean:
		$(GO) clean

test:
		$(GO) test .

check:
		@$(GO) list -f '{{join .Deps "\n"}}' | xargs $(GO) list -f '{{if not .Standard}}{{.ImportPath}} {{.Dir}}{{end}}' | column -t
