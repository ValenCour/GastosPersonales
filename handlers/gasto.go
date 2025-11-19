package handlers

import (
	sqlc "Tp3/db/generated"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GastosHandler(w http.ResponseWriter, r *http.Request) {
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
		err := r.ParseForm()
		if err != nil {
			fmt.Println("Error al decodificar usuario", err)
		}

		monto := r.FormValue("monto")
		medio_de_pago := r.FormValue("medio_de_pago")
		fechaStr := r.FormValue("fecha")
		categoria := r.FormValue("categoria")
		idStr := r.FormValue("id_usuario")

		fecha, err := time.Parse("2006-01-02T15:04", fechaStr)
		if err != nil {
			fmt.Println("Error parseando fecha:", err)
			return
		}
		if idStr != "" && idStr != "-1" {
			if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				if gastoValido(monto, medio_de_pago, fecha, sqlc.CategoriaGasto(categoria)) {
					gasto, err_create := queries.CreateGasto(r.Context(), sqlc.CreateGastoParams{
						IDUsuario:   int32(id),
						Monto:       monto,
						MedioDePago: medio_de_pago,
						Fecha:       fecha,
						Categoria:   sqlc.CategoriaGasto(categoria),
					})
					if err_create != nil {
						fmt.Println("Error al crear gasto", err_create)
					} else {
						fmt.Println("Gasto creado: ", gasto)
					}
				}
			} else {
				fmt.Println("Error al parsear id", err)
			}
			http.Redirect(w, r, "/?id_usuario="+idStr, http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func GastosIdHandler(w http.ResponseWriter, r *http.Request) {
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

func GastosDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	idUsuario := r.FormValue("id_usuario")
	id := strings.TrimPrefix(r.URL.Path, "/gastos/delete/")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error al convertir id", err)
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	err = queries.DeleteGasto(r.Context(), int32(id_int))
	if err != nil {
		fmt.Println("Error al eliminar gasto:", err)
		http.Error(w, "No se pudo eliminar el gasto", http.StatusInternalServerError)
		return
	} else {
		fmt.Printf("Gasto con id = %d eliminado \n", id_int)
	}
	http.Redirect(w, r, "/?id_usuario="+idUsuario, http.StatusSeeOther)
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
