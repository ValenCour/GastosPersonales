# Nombre para los servicios de Docker Compose
COMPOSE_PROJECT_NAME=tp3

# Target por defecto
default: test	

test: kill run script finish

script: 
	@bash -c "./prueba.sh"

up: 
	@docker compose up -d
	@echo "--- Esperando a que la base de datos esté lista ---"
	@until docker compose exec db pg_isready -U postgres -h localhost; do sleep 1; done
	@echo "--- Base de datos lista ---"

down: 
	@docker compose down

run: build up
	@./bin/app &
	@echo "--- Servidor iniciado ---"

generate: 
	@cd db/ && sqlc generate && cd ..
	@templ generate

kill:
	@pkill -f "./bin/app"

kill-hard:
	@bash -c "sudo lsof -ti :8080 | xargs -r sudo kill -9"

build: generate
	@go build -o ./bin/app .

clean: 
	@rm -rf ./bin ./tmp

finish: down clean kill
	
# Le dice a Make que estos no son archivos
.PHONY: default up down generate test run waitdocker kill kill-hard build clean finish