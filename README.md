Aplicación web en Go para llevar un registro de gastos personales.

`db` contiene todo lo relacionado con la definición y acceso a datos:

schema.sql: Contiene la definición de la tabla `usuarios`, `gastos` y el tipo `categoria_gasto`.

queries.sql: Incluye las consultas CRUD tanto para `usuarios` como para `gastos`, con anotaciones de `sqlc`.

sqlc.yaml Configuración necesaria para que `sqlc` lea el esquema y consultas, y genere el código Go en la carpeta `generated/`.

/generated: Contiene los archivos Go generados automáticamente.

## Ejecución

Abrir una terminal en el directorio raíz del proyecto y ejecutar en el siguiente orden:

docker run --name gastos-db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=gastos_db -p 5432:5432 -d docker.io/postgres

docker cp db/schema.sql gastos-db:/schema.sql

docker exec -it gastos-db psql -U postgres -d gastos_db -f /schema.sql

docker exec -it gastos-db psql -U postgres -d gastos_db

\q

cd db/

sqlc generate

cd ..

go run main.go

Esto permite comprobar el correcto funcionamiento de las queries generadas.

