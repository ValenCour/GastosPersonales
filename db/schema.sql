CREATE TYPE categoria_gasto AS ENUM ('Comida', 'Transporte', 'Entretenimiento', 'Servicios', 'Otros');

CREATE TABLE usuarios (
    id_usuario SERIAL PRIMARY KEY,
    nombre_usuario VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    contraseña VARCHAR(50) NOT NULL
);

CREATE TABLE gastos (
    id_gasto SERIAL PRIMARY KEY,
    id_usuario INT NOT NULL,
    monto DECIMAL(10,2) NOT NULL,
    medio_de_pago VARCHAR(50) NOT NULL,
    fecha TIMESTAMP NOT NULL,
    categoria categoria_gasto NOT NULL,
    FOREIGN KEY (id_usuario) REFERENCES usuarios(id_usuario)
);