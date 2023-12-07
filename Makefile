vet:
	go vet ./...
	shadow ./...

lint:
	golangci-lint run ./...

#run: vet, lint, migrate-up
#	go run ./cmd/app/main.go

clean:
	rm -rf ./build/* || true

build: clean
	go build -o ./build/app ./cmd/app/main.go

run:
	docker-compose up --remove-orphans --build app

migrate:
	migrate create -ext sql -dir migrations/ ${NAME}

migrate-up:
	migrate -source file://migrations -database postgres://root:secret@localhost:5432/backend?sslmode=disable up

migrate-down:
	migrate -source file://migrations -database postgres://root:secret@localhost:5432/backend?sslmode=disable down 1

migrate-force:
	migrate -source file://migrations -database postgres://root:secret@localhost:5432/backend?sslmode=disable force ${V}

swag:
	swag init -g ./cmd/app/main.go -o ./docs

.PHONY: vet, lint, migrate, migrate-up, migrate-down, migrate-force, swag,
