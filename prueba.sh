#!/bin/bash

# Este script prueba el flujo completo de la API de gastos.

echo "### Iniciando pruebas de la API ###"

curl -X POST http://localhost:8080/usuarios \
-H "Content-Type: application/json" \
-d '{
    "nombre_usuario": "valentino",
    "email": "valentino@example.com",
    "contraseña": "password123"
}'

curl -X POST http://localhost:8080/gastos \
-H "Content-Type: application/json" \
-d '{
    "id_usuario": 1,
    "monto": "1250.75",
    "medio_de_pago": "Tarjeta de Crédito",
    "fecha": "2025-10-27T10:30:00Z",
    "categoria": "Entretenimiento"
}'