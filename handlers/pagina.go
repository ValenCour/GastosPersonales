package handlers

import (
	sqlc "Tp3/db/generated"
	"Tp3/views"
	"fmt"
	"net/http"
	"strconv"
)

func PaginaHandler(w http.ResponseWriter, r *http.Request) {
	usuarios, err := queries.ListUsuarios(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener usuarios", http.StatusInternalServerError)
		return
	}

	var gastos []sqlc.Gasto
	var selectedID int32 = -1

	idStr := r.URL.Query().Get("id_usuario")

	if idStr != "" && idStr != "-1" {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			selectedID = int32(id)
			gastos, err = queries.ListGastosId(r.Context(), int32(id))
			if err != nil {
				fmt.Println("Error trayendo gastos:", err)
			}
		}
	}

	component := views.Estructura(usuarios, gastos, selectedID)

	component.Render(r.Context(), w)
}
