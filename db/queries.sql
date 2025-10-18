-- name: CreateUsuario :one
INSERT INTO usuarios (nombre_usuario, email, contraseña)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUsuario :one
SELECT * FROM usuarios
WHERE id_usuario = $1;

-- name: ListUsuarios :many
SELECT * FROM usuarios;

-- name: UpdateUsuario :one
UPDATE usuarios
SET nombre_usuario = $2,
    email = $3,
    contraseña = $4
WHERE id_usuario = $1
RETURNING *;

-- name: DeleteUsuario :exec
DELETE FROM usuarios
WHERE id_usuario = $1;

-- name: CreateGasto :one
INSERT INTO gastos (id_usuario, monto, medio_de_pago, fecha, categoria)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetGasto :one
SELECT * FROM gastos
WHERE id_gasto = $1;

-- name: ListGastos :many
SELECT * FROM gastos;

-- name: UpdateGasto :one
UPDATE gastos
SET monto = $2,
    medio_de_pago = $3,
    fecha = $4,
    categoria = $5
WHERE id_gasto = $1
RETURNING *;

-- name: DeleteGasto :exec
DELETE FROM gastos
WHERE id_gasto = $1;