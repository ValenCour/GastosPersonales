# Nombre para los servicios de Docker Compose
COMPOSE_PROJECT_NAME=tp3

# Target por defecto
default: test	

# Target principal para correr los tests de forma aislada
test:
	@bash -c "sudo lsof -ti :8080 | xargs -r sudo kill -9"
	@cd db/ && sqlc generate && cd ..

	@gnome-terminal --title="Servidor" -- bash -c "\
		docker compose up -d ; \
		go run main.go & \
		sleep 4 ; \
		read -p 'Presiona ENTER para terminar el servidor...'; \
		docker compose down ; \
		@bash -c "sudo lsof -ti :8080 | xargs -r sudo kill -9" ; \
		exec bash"

	@sleep 3

	@gnome-terminal --title="Cliente" -- bash -c "\
		./prueba.sh; \
		exec bash"

up: 
	@docker compose up -d

down: 
	@docker compose down
	@bash -c "sudo lsof -ti :8080 | xargs -r sudo kill -9"

run: up
	@go run main.go &

generate: 
	@cd db/ && sqlc generate && cd ..

# Le dice a Make que estos no son archivos
.PHONY: default build up down sqlc test run clean