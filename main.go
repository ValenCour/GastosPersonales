package main

import (
	sqlc "Tp3/db/generated"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

var queries *sqlc.Queries

func main() {
	var err_db error
	db, err_db = sql.Open("postgres", os.Getenv("DB_SOURCE"))
	if err_db != nil {
		log.Fatal("Error abriendo conexión:", err_db)
	} else {
		queries = sqlc.New(db)
		defer db.Close()
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/usuarios", usuariosHandler)

	mux.HandleFunc("/usuarios/", usuariosIdHandler)

	mux.HandleFunc("/gastos", gastosHandler)

	mux.HandleFunc("/gastos/", gastosIdHandler)

	log_mux := loggingMiddleware(mux)

	err := http.ListenAndServe(":8080", log_mux)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}

func usuariosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		usuarios, err := queries.ListUsuarios(r.Context())
		if err != nil {
			fmt.Println("Error al obtener usuarios de la base", err)
		}

		err = json.NewEncoder(w).Encode(usuarios)
		if err != nil {
			fmt.Println("Error al codificar usuarios", err)
		} else {
			fmt.Println("Lista de usuarios:", usuarios)
		}
	case http.MethodPost:
		var nuevo_usuario sqlc.CreateUsuarioParams
		err := json.NewDecoder(r.Body).Decode(&nuevo_usuario)
		if err != nil {
			fmt.Println("Error al decodificar usuario", err)
		}

		w.Header().Set("Content-Type", "application/json")

		if usuarioValido(nuevo_usuario.NombreUsuario, nuevo_usuario.Email, nuevo_usuario.Contraseña) {
			usuario, err_create := queries.CreateUsuario(r.Context(), sqlc.CreateUsuarioParams{
				NombreUsuario: nuevo_usuario.NombreUsuario,
				Email:         nuevo_usuario.Email,
				Contraseña:    nuevo_usuario.Contraseña,
			})
			if err_create != nil {
				fmt.Println("Error al crear usuario", err_create)
			} else {
				w.WriteHeader(http.StatusCreated)
				fmt.Println("Usuario creado: ", usuario)
			}
		}
	}
}

func usuariosIdHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/usuarios/")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error al convertir id", err)
	}

	usuario_encontrado, err_get := queries.GetUsuario(r.Context(), int32(id_int))
	if err_get != nil {
		fmt.Println("Error al obtener usuario de la base", err)
		http.NotFound(w, r)
	} else {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")

			err = json.NewEncoder(w).Encode(usuario_encontrado)
			if err != nil {
				fmt.Println("Error al codificar usuario", err)
			} else {
				fmt.Println("Usuario encontrado: ", usuario_encontrado)
			}
		case http.MethodPut:
			var usuario_mod sqlc.UpdateUsuarioParams
			err := json.NewDecoder(r.Body).Decode(&usuario_mod)
			if err != nil {
				fmt.Println("Error al decodificar gasto", err)
			}

			if usuarioValido(usuario_mod.NombreUsuario, usuario_mod.Email, usuario_mod.Contraseña) {
				usuario, err := queries.UpdateUsuario(r.Context(), sqlc.UpdateUsuarioParams{
					IDUsuario:     int32(id_int),
					NombreUsuario: usuario_mod.NombreUsuario,
					Email:         usuario_mod.Email,
					Contraseña:    usuario_mod.Contraseña,
				})
				if err != nil {
					fmt.Println("Error al actualizar usuario:", err)
				} else {
					fmt.Println("Usuario modificado: ", usuario)
				}
			}
		case http.MethodDelete:
			err := queries.DeleteUsuario(r.Context(), int32(id_int))
			if err != nil {
				fmt.Println("Error al eliminar usuario:", err)
			} else {
				w.WriteHeader(http.StatusNoContent)
				fmt.Printf("Usuario con id = %d eliminado \n", id_int)
			}
		}
	}
}

func gastosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")

		gastos, err := queries.ListGastos(r.Context())
		if err != nil {
			fmt.Println("Error al obtener gastos de la base", err)
		}

		err = json.NewEncoder(w).Encode(gastos)
		if err != nil {
			fmt.Println("Error al codificar gastos", err)
		} else {
			fmt.Println("Lista de gastos:", gastos)
		}
	case http.MethodPost:
		var nuevo_gasto sqlc.CreateGastoParams
		err := json.NewDecoder(r.Body).Decode(&nuevo_gasto)
		if err != nil {
			fmt.Println("Error al decodificar gasto", err)
		}

		if gastoValido(nuevo_gasto.Monto, nuevo_gasto.MedioDePago, nuevo_gasto.Fecha, nuevo_gasto.Categoria) {
			w.Header().Set("Content-Type", "application/json")

			gasto, err_create := queries.CreateGasto(r.Context(), sqlc.CreateGastoParams{
				IDUsuario:   nuevo_gasto.IDUsuario,
				Monto:       nuevo_gasto.Monto,
				MedioDePago: nuevo_gasto.MedioDePago,
				Fecha:       nuevo_gasto.Fecha,
				Categoria:   nuevo_gasto.Categoria,
			})
			if err_create != nil {
				fmt.Println("Error al crear gasto", err_create)
			} else {
				w.WriteHeader(http.StatusCreated)
				fmt.Println("Gasto creado: ", gasto)
			}
		}
	}
}

func gastosIdHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/gastos/")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error al convertir id", err)
	}

	gasto_encontrado, err_get := queries.GetGasto(r.Context(), int32(id_int))
	if err_get != nil {
		fmt.Println("Error al obtener gasto de la base", err)
		http.NotFound(w, r)
	} else {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(gasto_encontrado)
			if err != nil {
				fmt.Println("Error al codificar gasto", err)
			} else {
				fmt.Println("Gasto encontrado: ", gasto_encontrado)
			}
		case http.MethodPut:
			var gasto_mod sqlc.UpdateGastoParams
			err := json.NewDecoder(r.Body).Decode(&gasto_mod)

			if err != nil {
				fmt.Println("Error al decodificar gasto", err)
			}

			if gastoValido(gasto_mod.Monto, gasto_mod.MedioDePago, gasto_mod.Fecha, gasto_mod.Categoria) {
				gasto, err := queries.UpdateGasto(r.Context(), sqlc.UpdateGastoParams{
					IDGasto:     int32(id_int),
					Monto:       gasto_mod.Monto,
					MedioDePago: gasto_mod.MedioDePago,
					Fecha:       gasto_mod.Fecha,
					Categoria:   gasto_mod.Categoria,
				})

				if err != nil {
					fmt.Println("Error al actualizar gasto:", err)
				} else {
					fmt.Println("Gasto modificado: ", gasto)
				}
			}
		case http.MethodDelete:
			err := queries.DeleteGasto(r.Context(), int32(id_int))
			if err != nil {
				fmt.Println("Error al eliminar gasto:", err)
			} else {
				w.WriteHeader(http.StatusNoContent)
				fmt.Printf("Gasto con id = %d eliminado \n", id_int)
			}
		}
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Petición recibida: %s %s \n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func usuarioValido(nombre, email, contraseña string) bool {
	if nombre != "" && email != "" && contraseña != "" {
		fmt.Println("Usuario válido")
		return true
	} else {
		fmt.Println("Usuario inválido")
		return false
	}
}

func gastoValido(monto string, medio_pago string, fecha time.Time, categoria sqlc.CategoriaGasto) bool {
	if monto != "" && medio_pago != "" && !fecha.IsZero() && categoria != "" {
		fmt.Println("Gasto válido")
		return true
	} else {
		fmt.Println("Gasto inválido")
		return false
	}
}
