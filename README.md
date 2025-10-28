# Proyecto: Control de Gastos Personales

Aplicación web en Go para llevar un registro de gastos personales.

# Estructura
- `db/`: Contiene todo lo relacionado con la base de datos. 
- `handlers/`: Lógica de los Endpoints de la API.
- `public/`: Archivos del Frontend.
- `main.go`: Punto de entrada de la aplicación: configura y arranca el servidor.
- `Makefile`: Automatización de tareas.

# Ejecución



# Comandos Disponibles (Makefile)

Se puede usar `make` para ejecutar diversas tareas:
- `make run`: Construye y ejecuta la aplicación completa con su base de datos.
- `make test`: Inicia servidor, ejecuta la secuencia de pruebas definida en `prueba.sh` y luego limpia el entorno.
- `make script`: Ejecuta el script de prueba `prueba.sh`.
- `make down`: Detiene y elimina los contenedores de Docker.
- `make build`: Compila el código fuente de Go y genera el ejecutable en `./bin/app`.
- `make generate`: Ejecuta `sqlc` para generar el código Go a partir de los archivos `.sql`.
- `make kill`: Detiene el proceso del servidor si se está ejecutando localmente.
- `make clean`: Elimina los directorios de compilación.
- `make finish`: Ejecuta una limpieza completa del entorno (`down`, `clean`, `kill`).