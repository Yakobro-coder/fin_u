.PHONY: run
run:
	docker-compose --env-file ./env/dev.env up -d

.PHONY: local_run
local_run:
	docker run --name db_local --env-file ./env/dev.env -d -p 5432:5432 postgres
	go run cmd/main.go

.PHONY: migrate
DB_URI=postgres://user:PasswordIsPassword@localhost:5432/material_school
migrate: install-tools
	goose -dir ./migrations postgres "$(DB_URI)" up -v

.PHONY: install-tools
install-tools:
	go install github.com/pressly/goose/v3/cmd/goosegit init@v3.20.0