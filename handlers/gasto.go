package handlers

import (
	sqlc "Tp3/db/generated"
	"Tp3/views"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GastosHandler(w http.ResponseWriter, r *http.Request) {
	id_usuario := r.Context().Value("userID").(int32)
	switch r.Method {
	case http.MethodGet:
		gastos, err := queries.ListGastosId(r.Context(), id_usuario)
		if err != nil {
			http.Error(w, "Error al obtener gastos", 500)
			return
		}

		sortColumn := r.URL.Query().Get("sort")
		sortOrder := r.URL.Query().Get("order")

		if sortOrder == "" {
			sortOrder = "asc"
		}

		sort.Slice(gastos, func(i, j int) bool {
			g1, g2 := gastos[i], gastos[j]

			menor := false
			switch sortColumn {
			case "monto":
				m1, _ := strconv.ParseFloat(g1.Monto, 64)
				m2, _ := strconv.ParseFloat(g2.Monto, 64)
				menor = m1 < m2
			case "fecha":
				menor = g1.Fecha.Before(g2.Fecha)
			case "categoria":
				menor = string(g1.Categoria) < string(g2.Categoria)
			default:
				return g1.IDGasto < g2.IDGasto
			}

			if sortOrder == "desc" {
				return !menor
			}
			return menor
		})
		if r.Header.Get("HX-Request") == "true" {
			views.GastosList(gastos, id_usuario, sortColumn, sortOrder).Render(r.Context(), w)
		} else {
			usuario, _ := queries.GetUsuario(r.Context(), id_usuario)
			views.Estructura(usuario, gastos, id_usuario).Render(r.Context(), w)
		}
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error procesando formulario", 400)
			return
		}
		monto := r.FormValue("monto")
		medio_de_pago := r.FormValue("medio_de_pago")
		fechaStr := r.FormValue("fecha")
		categoria := r.FormValue("categoria")

		fecha, err := time.Parse("2006-01-02T15:04", fechaStr)
		if err != nil {
			http.Error(w, "Fecha inválida", 400)
			return
		}

		if !gastoValido(monto, medio_de_pago, fecha, sqlc.CategoriaGasto(categoria)) {
			http.Error(w, "Datos inválidos", 400)
			return
		}

		_, err = queries.CreateGasto(r.Context(), sqlc.CreateGastoParams{
			IDUsuario:   id_usuario,
			Monto:       monto,
			MedioDePago: medio_de_pago,
			Fecha:       fecha,
			Categoria:   sqlc.CategoriaGasto(categoria),
		})

		if err != nil {
			http.Error(w, "Error al insertar gasto", 500)
			return
		}

		gastos, err := queries.ListGastosId(r.Context(), id_usuario)
		if err != nil {
			http.Error(w, "Error obteniendo lista", 500)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		views.GastosList(gastos, int32(id_usuario), "", "").Render(r.Context(), w)
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
				http.Error(w, "Error al eliminar el gasto", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Printf("Gasto con id = %d eliminado \n", id_int)
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
