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

var db, err_db = sql.Open("postgres", os.Getenv("DB_SOURCE"))

var queries = sqlc.New(db)

func main() {
	if err_db != nil {
		log.Fatal("Error abriendo conexión:", err_db)
	}
	defer db.Close()

	http.HandleFunc("/usuarios", usuariosHandler)

	http.HandleFunc("/usuarios/", usuariosIdHandler)

	http.HandleFunc("/gastos", gastosHandler)

	http.HandleFunc("/gastos/", gastosIdHandler)

	err := http.ListenAndServe(":8080", nil)
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
		}
	case http.MethodPost:
		var nuevo_usuario sqlc.CreateUsuarioParams
		err := json.NewDecoder(r.Body).Decode(&nuevo_usuario)
		if err != nil {
			fmt.Println("Error al decodificar usuario", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err_create := queries.CreateUsuario(r.Context(), sqlc.CreateUsuarioParams{
			NombreUsuario: nuevo_usuario.NombreUsuario,
			Email:         nuevo_usuario.Email,
			Contraseña:    nuevo_usuario.Contraseña,
		})
		if err_create != nil {
			fmt.Println("Error al crear usuario", err_create)
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
			}
		case http.MethodPut:
			var usuario_mod sqlc.UpdateUsuarioParams
			err := json.NewDecoder(r.Body).Decode(&usuario_mod)

			if err != nil {
				fmt.Println("Error al decodificar gasto", err)
			}

			_, err = queries.UpdateUsuario(r.Context(), sqlc.UpdateUsuarioParams{
				IDUsuario:     int32(id_int),
				NombreUsuario: usuario_mod.NombreUsuario,
				Email:         usuario_mod.Email,
				Contraseña:    usuario_mod.Contraseña,
			})

			if err != nil {
				fmt.Println("Error al actualizar usuario:", err)
			}
		case http.MethodDelete:
			err := queries.DeleteUsuario(r.Context(), int32(id_int))
			if err != nil {
				fmt.Println("Error al eliminar usuario:", err)
			} else {
				w.WriteHeader(http.StatusNoContent)
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
		}
	case http.MethodPost:
		var nuevo_gasto sqlc.CreateGastoParams
		err := json.NewDecoder(r.Body).Decode(&nuevo_gasto)
		if err != nil {
			fmt.Println("Error al decodificar gasto", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err_create := queries.CreateGasto(r.Context(), sqlc.CreateGastoParams{
			IDUsuario:   nuevo_gasto.IDUsuario,
			Monto:       nuevo_gasto.Monto,
			MedioDePago: nuevo_gasto.MedioDePago,
			Fecha:       time.Now(),
			Categoria:   nuevo_gasto.Categoria,
		})
		if err_create != nil {
			fmt.Println("Error al crear gasto", err_create)
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
		fmt.Println("Error al obtener gastos de la base", err)
		http.NotFound(w, r)
	} else {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(gasto_encontrado)
			if err != nil {
				fmt.Println("Error al codificar gasto", err)
			}
		case http.MethodPut:
			var gasto_mod sqlc.UpdateGastoParams
			err := json.NewDecoder(r.Body).Decode(&gasto_mod)

			if err != nil {
				fmt.Println("Error al decodificar gasto", err)
			}

			_, err = queries.UpdateGasto(r.Context(), sqlc.UpdateGastoParams{
				IDGasto:     int32(id_int),
				Monto:       gasto_mod.Monto,
				MedioDePago: gasto_mod.MedioDePago,
				Fecha:       gasto_mod.Fecha,
				Categoria:   gasto_mod.Categoria,
			})

			if err != nil {
				fmt.Println("Error al actualizar gasto:", err)
			}
		case http.MethodDelete:
			err := queries.DeleteGasto(r.Context(), int32(id_int))
			if err != nil {
				fmt.Println("Error al eliminar gasto:", err)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		}
	}
}
