package main

import (
	sqlc "Tp3/db/generated"
	handlers "Tp3/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB

var queries *sqlc.Queries

func main() {
	var err_db error
	db, err_db = sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/gastos_db?sslmode=disable")

	if err_db != nil {
		log.Fatal("Error abriendo conexión:", err_db)
	} else {
		queries = sqlc.New(db)
		handlers.SetQueries(queries)
		defer db.Close()
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.PaginaHandler)

	mux.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir("views"))))

	mux.HandleFunc("/usuarios", handlers.UsuariosHandler)

	mux.HandleFunc("/usuarios/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/gastos") {
			handlers.GastosPorUsuarioHandler(w, r)
		} else {
			handlers.UsuariosIdHandler(w, r)
		}
	})

	mux.HandleFunc("/gastos", handlers.GastosHandler)

	mux.HandleFunc("/gastos/", handlers.GastosIdHandler)

	mux.HandleFunc("/gastos/delete/", handlers.GastosDeleteHandler)

	log_mux := loggingMiddleware(mux)

	err := http.ListenAndServe(":8080", log_mux)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Petición recibida: %s %s \n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
