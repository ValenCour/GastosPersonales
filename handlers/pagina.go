package handlers

import (
	sqlc "Tp3/db/generated"
	"Tp3/views"
	"fmt"
	"net/http"
)

func PaginaHandler(w http.ResponseWriter, r *http.Request) {
	var gastos []sqlc.Gasto
	id_usuario := r.Context().Value("userID").(int32)

	usuario, err := queries.GetUsuario(r.Context(), id_usuario)
	if err != nil {
		http.Error(w, "Error al obtener usuarios", http.StatusInternalServerError)
		return
	}

	gastos, err = queries.ListGastosId(r.Context(), id_usuario)
	if err != nil {
		fmt.Println("Error trayendo gastos:", err)
	}

	component := views.Estructura(usuario, gastos, id_usuario)

	component.Render(r.Context(), w)
}
