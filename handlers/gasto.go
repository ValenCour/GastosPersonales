package handlers

import (
	sqlc "Tp3/db/generated"
	views "Tp3/views"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
)

func GastosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")

		gastos, err := queries.ListGastos(r.Context())
		if err != nil {
			fmt.Println("Error al obtener gastos de la base", err)
		}
		templ.Handler(views.GastosPage(gastos)).ServeHTTP(w, r)

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

func gastoValido(monto string, medio_pago string, fecha time.Time, categoria sqlc.CategoriaGasto) bool {
	if monto != "" && medio_pago != "" && !fecha.IsZero() && categoria != "" {
		fmt.Println("Gasto válido")
		return true
	} else {
		fmt.Println("Gasto inválido")
		return false
	}
}
