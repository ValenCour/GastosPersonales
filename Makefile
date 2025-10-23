# Nombre para los servicios de Docker Compose
COMPOSE_PROJECT_NAME=tp3

# Target por defecto
default: test	

test: kill generate run wait-docker script down kill

script: 
	@bash -c "./prueba.sh"

up: 
	@docker compose up -d

down: 
	@docker compose down
	@bash -c "sudo lsof -ti :8080 | xargs -r sudo kill -9"

wait-docker:
	@until docker compose exec db pg_isready -U postgres -h localhost; do sleep 1; done

run: up
	@go run main.go &
	@echo "Iniciando servidor"
	@sleep 5

generate: 
	@cd db/ && sqlc generate && cd ..

kill:
	@bash -c "sudo lsof -ti :8080 | xargs -r sudo kill -9"

# Le dice a Make que estos no son archivos
.PHONY: default up down generate test run wait-docker wait-server test1 test2