# Nombre para los servicios de Docker Compose
COMPOSE_PROJECT_NAME=tp3

# Target por defecto
default: test	

# Target principal para correr los tests de forma aislada
test:
	@cd db/ && sqlc generate && cd ..
	@docker compose up -d

	@gnome-terminal --title="Servidor Go (Logs)" -- bash -c "\
		export DB_SOURCE='postgres://postgres:postgres@localhost:5432/gastos_db?sslmode=disable'; \
		go run main.go & \
		sleep 7 ; \
		read -p 'Presiona ENTER para terminar el servidor...'; \
		docker compose down ; \
		pkill -f "main" || true ; \
		exec bash"

	@sleep 5

	@gnome-terminal --title="Cliente" -- bash -c "\
		./prueba.sh; \
		exec bash"

# Le dice a Make que estos no son archivos
.PHONY: default build up down sqlc test run clean
