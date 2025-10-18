# Nombre para los servicios de Docker Compose
COMPOSE_PROJECT_NAME=tp3

# Target por defecto, el que se ejecuta si solo escribes "make"
default: test	

# Target principal para correr los tests de forma aislada
test:
	@cd db/ && sqlc generate && cd ..
	@docker compose up -d
	@DB_SOURCE="postgres://postgres:postgres@localhost:5432/gastos_db?sslmode=disable" go run main.go & sleep 2
	@./prueba.sh
	@docker compose down
	@pkill -f "main" || true

# Le dice a Make que estos no son archivos, sino comandos
.PHONY: default build up down sqlc test run clean