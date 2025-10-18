#!/bin/bash

# Este script prueba el flujo completo de la API de gastos.

echo "### Inicio de pruebas ###"

echo "=== POSTs validos en /usuarios ==="

curl -X POST http://localhost:8080/usuarios \
-H "Content-Type: application/json" \
-d '{
    "nombre_usuario": "usuario1",
    "email": "1@example.com",
    "contraseña": "clave123"
}'

curl -X POST http://localhost:8080/usuarios \
-H "Content-Type: application/json" \
-d '{
    "nombre_usuario": "usuario2",
    "email": "2@example.com",
    "contraseña": "clave456"
}'

echo "=== POST invalido en /usuarios ==="

curl -X POST http://localhost:8080/usuarios \
-H "Content-Type: application/json" \
-d '{
    "nombre_usuario": "usuario3",
    "email": "3@example.com",
    "contraseña": ""
}'

echo "=== GET en /usuarios ==="

curl -X GET http://localhost:8080/usuarios 

echo "=== GET válido en /usuarios/1 ==="

curl -X GET http://localhost:8080/usuarios/1

echo "=== GET inválido en /usuarios/10 ==="

curl -X GET http://localhost:8080/usuarios/10

echo "=== PUT válido en /usuarios/1 ==="

curl -X PUT http://localhost:8080/usuarios/1 \
   -H "Content-Type: application/json" \
   -d '{
    "nombre_usuario": "usuario1_actualizado",
    "email": "1_actualizado@example.com",
    "contraseña": "clave123_actualizado"
}'

echo "=== PUT inválido en /usuarios/2 ==="

curl -X PUT http://localhost:8080/usuarios/2 \
   -H "Content-Type: application/json" \
   -d '{
    "nombre_usuario": "",
    "email": "2_actualizado@example.com",
    "contraseña": "clave456_actualizado"
}'

echo "=== DELETE válido en /usuarios/2 ==="

curl -X DELETE http://localhost:8080/usuarios/2

echo "=== GET en /usuarios ==="

curl -X GET http://localhost:8080/usuarios

echo "=== POSTs validos en /gastos ==="

curl -X POST http://localhost:8080/gastos \
-H "Content-Type: application/json" \
-d '{
    "id_usuario": 1,
    "monto": "500.00",
    "medio_de_pago": "Efectivo",
    "fecha": "2025-01-01T12:00:00Z",
    "categoria": "Transporte"
}'

curl -X POST http://localhost:8080/gastos \
-H "Content-Type: application/json" \
-d '{
    "id_usuario": 2,
    "monto": "100.00",
    "medio_de_pago": "Transferencia",
    "fecha": "2025-01-01T12:00:00Z",
    "categoria": "Servicios"
}'

echo "=== POST invalido en /gastos ==="

curl -X POST http://localhost:8080/gastos \
-H "Content-Type: application/json" \
-d '{
    "id_usuario": 1,
    "monto": "",
    "medio_de_pago": "Tarjeta de Débito",
    "fecha": "",
    "categoria": "Otros"
}'

echo "=== GET en /gastos ==="

curl -X GET http://localhost:8080/gastos 

echo "=== GET válido en /gastos/1 ==="

curl -X GET http://localhost:8080/gastos/1

echo "=== GET inválido en /gastos/10 ==="

curl -X GET http://localhost:8080/gastos/10

echo "=== PUT válido en /gastos/1 ==="

curl -X PUT http://localhost:8080/gastos/1 \
-H "Content-Type: application/json" \
-d '{
    "monto": "125.50",
    "medio_de_pago": "Tarjeta de Débito",
    "fecha": "2025-10-28T20:00:00Z",
    "categoria": "Transporte"
}'

echo "=== PUT inválido en /gastos/2 ==="

curl -X PUT http://localhost:8080/gastos/2 \
-H "Content-Type: application/json" \
-d '{
    "monto": "4500.00",
    "medio_de_pago": "Tarjeta de Crédito",
    "fecha": "2025-10-29T11:30:00Z",
    "categoria": ""
}'

echo "=== DELETE válido en /gastos/2 ==="

curl -X DELETE http://localhost:8080/gastos/2

echo "=== GET en /gastos ==="

curl -X GET http://localhost:8080/gastos 

echo "### Fin de pruebas ###"