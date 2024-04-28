.PHONY: compose-up
compose-up:
	docker compose up --build && docker compose logs --follow

.PHONY: compose-down
compose-down:
	docker compose down --remove-orphans

.PHONY: migrate-create
migrate-create:
	migrate create -ext sql -seq -dir migrations 'init_schema'

.PHONY: migrate-up
migrate-up:
	migrate -path migrations -database 'postgres://postgres:password@localhost:5432/messages?sslmode=disable' up
