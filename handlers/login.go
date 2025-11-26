package handlers

import (
	sqlc "Tp3/db/generated"
	"Tp3/views"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type sesion struct {
	id_usuario int32
	expiracion time.Time
}

var sesiones = map[string]sesion{}

func AutenticacionHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := verificarSesion(r); err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	views.PaginaLogin("").Render(r.Context(), w)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := queries.GetUsuarioByEmail(r.Context(), email)
	if err != nil {
		views.PaginaLogin("Email inválido").Render(r.Context(), w)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Contraseña), []byte(password))
	if err != nil {
		views.PaginaLogin("Contraseña inválida").Render(r.Context(), w)
		return
	}

	crearSesion(w, user.IDUsuario)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RegistroHandler(w http.ResponseWriter, r *http.Request) {
	nombre := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if usuarioValido(nombre, email, password) {
		usuario, err_create := queries.CreateUsuario(r.Context(), sqlc.CreateUsuarioParams{
			NombreUsuario: nombre,
			Email:         email,
			Contraseña:    string(hashed),
		})
		if err_create != nil {
			fmt.Println("Error al crear usuario", err_create)
			views.PaginaLogin("Error: El email ya está registrado").Render(r.Context(), w)
			return
		} else {
			fmt.Println("Usuario creado: ", usuario)
			crearSesion(w, usuario.IDUsuario)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("session_token"); err == nil {
		delete(sesiones, c.Value)
	}
	http.SetCookie(w, &http.Cookie{Name: "session_token", Value: "", MaxAge: -1})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := verificarSesion(r)
		if err != nil {
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/login")
				w.WriteHeader(http.StatusOK)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		next(w, r.WithContext(ctx))
	}
}

func crearSesion(w http.ResponseWriter, id_usuario int32) {
	token := uuid.NewString()
	expiracion_sesion := time.Now().Add(300 * time.Second)
	sesiones[token] = sesion{
		id_usuario: id_usuario,
		expiracion: expiracion_sesion,
	}
	http.SetCookie(w, &http.Cookie{
		Name: "session_token", Value: token, Expires: expiracion_sesion, HttpOnly: true,
	})
}

func verificarSesion(r *http.Request) (int32, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0, err
	}

	token := cookie.Value
	sesion_usuario, exists := sesiones[token]

	if !exists {
		return 0, errors.New("sesión inválida")
	}

	if time.Now().After(sesion_usuario.expiracion) {
		delete(sesiones, token)
		return 0, errors.New("sesión expirada")
	}

	return sesion_usuario.id_usuario, nil
}
