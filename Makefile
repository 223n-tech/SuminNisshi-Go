.PHONY: lint tools.get
GOBIN = $(shell go env GOPATH)/bin

mod.install:
	go mod tidy

mod.install.adminlte:
	bash /usr/local/bin/setup-adminlte.sh

tools.get:
	go install github.com/mgechev/revive@latest

tools.doc:
	mkdir -p doc
	go run ./tools/doc-template-generator.go ./cmd/suiminnisshi ./doc/godoc/package-template.md ./doc/godoc/cmd.md
	go run ./tools/doc-template-generator.go ./internal/config ./doc/godoc/package-template.md ./doc/godoc/config.md
	go run ./tools/doc-template-generator.go ./internal/handler ./doc/godoc/package-template.md ./doc/godoc/handler.md
	go run ./tools/doc-template-generator.go ./internal/middleware ./doc/godoc/package-template.md ./doc/godoc/middleware.md
	go run ./tools/doc-template-generator.go ./internal/models ./doc/godoc/package-template.md ./doc/godoc/models.md
	go run ./tools/doc-template-generator.go ./internal/pdf ./doc/godoc/package-template.md ./doc/godoc/pdf.md
	go run ./tools/doc-template-generator.go ./internal/repository ./doc/godoc/package-template.md ./doc/godoc/repository.md
	go run ./tools/doc-template-generator.go ./internal/service ./doc/godoc/package-template.md ./doc/godoc/service.md
	go run ./tools/doc-template-generator.go ./internal/util ./doc/godoc/package-template.md ./doc/godoc/util.md
	go run ./tools/doc-template-generator.go ./tools ./doc/godoc/package-template.md ./doc/godoc/tools.md

revive:
	$(GOBIN)/revive -config revive.toml -formatter friendly ./...

vet:
	go vet ./...

lint: vet revive

run:
	go run cmd/suiminnisshi/main.go

version:
	go version
