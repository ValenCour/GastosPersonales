# Proyecto: Control de Gastos Personales

Aplicación web en Go para llevar un registro de gastos personales.

# Estructura
- `db/`: Contiene todo lo relacionado con la base de datos. 
- `handlers/`: Lógica de los Endpoints de la API.
- `views/`: Archivos .templ que definen la estructura HTML de las páginas.
- `main.go`: Punto de entrada de la aplicación: configura y arranca el servidor.
- `Makefile`: Automatización de tareas.

# Ejecución

1. **Levantar el servidor:**
    `make run`
    Este comando se encarga de:
    * Generar el código Go a partir de SQL.
    * Generar código Go mediante templ generate
    * Compilar la aplicación Go.
    * Iniciar la base de datos en un contenedor Docker.
    * Ejecuta el servidor de Go localmente.

2.  **Acceso a la aplicación:**
    Una vez que se muestre por consola el mensaje "--- Servidor iniciado ---", abrir el navegador en:[http://localhost:8080]

3.  **Limpieza del entorno**
    `make finish` permite limpiar completamente el entorno para detener todos los servicios y eliminar los archivos generados.

# Comandos Disponibles (Makefile)

Se puede usar `make` para ejecutar diversas tareas:
- `make run`: Construye y ejecuta la aplicación completa con su base de datos.
- `make down`: Detiene y elimina los contenedores de Docker.
- `make build`: Compila el código fuente de Go y genera el ejecutable en `./bin/app`.
- `make generate`: Ejecuta `sqlc` para generar el código Go a partir de los archivos `.sql`.
- `make templ`: Ejecuta `templ generate` para generar código Go a partir de los archivos `.templ`.
- `make kill`: Detiene el proceso del servidor si se está ejecutando localmente.
- `make clean`: Elimina los directorios de compilación.
- `make finish`: Ejecuta una limpieza completa del entorno (`down`, `clean`, `kill`).

# Endpoints de la API
La aplicación expone los siguientes endpoints:
- `GET /usuarios`: Lista todos los usuarios registrados.
- `POST /usuarios`: Crea un nuevo usuario.
- `GET /usuarios/{id}`: Obtiene la información de un usuario específico por su id.
- `PUT /usuarios/{id}`: Actualiza la información de un usuario existente.
- `DELETE /usuarios/{id}`: Elimina un usuario por su id.
- `GET /usuarios/{id}/gastos`: Lista todos los gastos de un usuario por su id.
- `GET /gastos`: Lista todos los gastos registrados.
- `POST /gastos`: Registra un nuevo gasto para un usuario.
- `GET /gastos/{id}`: Obtiene la información de un gasto específico por su id.
- `PUT /gastos/{id}`: Actualiza la información de un gasto existente.
- `DELETE /gastos/{id}`: Elimina un gasto por su id.
- `POST /gastos/delete/{id}`: Elimina un gasto por su id.