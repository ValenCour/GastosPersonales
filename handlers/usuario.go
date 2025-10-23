package handlers

import (
	sqlc "Tp3/db/generated"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func UsuariosHandler(w http.ResponseWriter, r *http.Request) {
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

func UsuariosIdHandler(w http.ResponseWriter, r *http.Request) {
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

func usuarioValido(nombre, email, contraseña string) bool {
	if nombre != "" && email != "" && contraseña != "" {
		fmt.Println("Usuario válido")
		return true
	} else {
		fmt.Println("Usuario inválido")
		return false
	}
}

var queries *sqlc.Queries

func SetQueries(q *sqlc.Queries) {
	queries = q
}
