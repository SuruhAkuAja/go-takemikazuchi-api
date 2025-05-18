# Gunakan shell bash
SHELL := /bin/bash

# Variabel
GO := go
MIGRATE_COMMAND := migrate
RUN_CMD := $(GO) run cmd/server/main.go
table = $(table)
name = $(name)
url=mysql://root:root@tcp(127.0.0.1:3306)/go_takemikazuchi_api
version=$(version)

# Command utama
.PHONY: run build clean test fmt vet lint wire gen inject migrate migration-up


# Menjalankan aplikasi
default: run

run:
	$(RUN_CMD)

# Build aplikasi
build:
	$(GO) build -o bin/app cmd/server/main.go

# Membersihkan file hasil build
clean:
	rm -rf bin/

# Menjalankan unit test
test:
	$(GO) test ./... -cover

# Format kode
dfmt:
	$(GO) fmt ./...

# Static analysis
vet:
	$(GO) vet ./...

# Linting menggunakan golangci-lint
lint:
	golangci-lint run ./...

# Menjalankan semua check sekaligus
check: fmt vet lint test

# Build injection
inject:
	wire gen ./cmd/injection/injector.go

migration-up:
	$(MIGRATE_COMMAND) -database "$(url)" -path ./migrations up

migration-down:
	$(MIGRATE_COMMAND) -database "$(url)" -path ./migrations down

migration-up-one:
	$(MIGRATE_COMMAND) -database "$(url)" -path ./migrations up 1

migration-down-one:
	$(MIGRATE_COMMAND) -database "$(url)" -path ./migrations down 1

migration-create:
	$(MIGRATE_COMMAND) create -ext sql -dir ./migrations $(name)

migration-force:
	$(MIGRATE_COMMAND) -database "$(url)" -path ./migrations force $(version)