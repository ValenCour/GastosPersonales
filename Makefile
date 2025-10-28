COMPOSE_PROJECT_NAME=tp3

default: test	

#Inicia servidor, ejecuta las pruebas y termina
test: kill run script finish

#Ejecuta el script de prueba ubicado en 'prueba.sh'
script: 
	@bash -c "./prueba.sh"

#Levanta la base de datos y espera hasta que esté lista antes de continuar
up: 
	@docker compose up -d
	@echo "--- Esperando a que la base de datos esté lista ---"
	@until docker compose exec db pg_isready -U postgres -h localhost; do sleep 1; done
	@echo "--- Base de datos lista ---"

#Detiene y elimina los contenedores creados
down: 
	@docker compose down

#Ejecuta el binario de la aplicación en segundo plano
run: build up
	@./bin/app &
	@echo "--- Servidor iniciado ---"

#Genera el código Go a partir de las consultas SQL usando sqlc.
generate: 
	@cd db/ && sqlc generate && cd ..

#Mata cualquier proceso que esté corriendo './bin/app'.
kill:
	@pkill -f "./bin/app"

#Busca cualquier proceso que esté usando el puerto 8080 y lo termina
#Usar si 'kill' falla
kill-hard:
	@bash -c "sudo lsof -ti :8080 | xargs -r sudo kill -9"

#Compila el código Go
build: generate
	@go build -o ./bin/app .

#Elimina los directorios de binarios y temporales
clean: 
	@rm -rf ./bin ./tmp

#Regla de finalización que: elimina los contenedores, borra los binarios y mata cualquier proceso del servidor
finish: down clean kill

.PHONY: default up down generate test run waitdocker kill kill-hard build clean finish