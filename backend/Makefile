DB_URL=postgres://bookshop:bookshop123@localhost:5432/bookshop?sslmode=disable
MIGRATIONS=./migrations

hot-reload:
	air -c .air.toml

all: hot-reload

migrate-new:
	migrate create -ext sql -dir $(MIGRATIONS) -seq $(name)

migrate-up:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" down

migrate-drop:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" drop -f

migrate-reset:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" drop -f
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" up

.PHONY: migrate-new migrate-up migrate-down migrate-drop migrate-reset hot-reload all

