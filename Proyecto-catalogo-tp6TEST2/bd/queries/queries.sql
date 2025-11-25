-- name: GetCliente :one
SELECT * FROM clientes WHERE id = $1;

-- name: ListClientes :many
SELECT * FROM clientes ORDER BY email;

-- name: GetArticulo :one
SELECT * FROM articulos WHERE id = $1;

-- name: ListArticulos :many
SELECT * FROM articulos ORDER BY titulo;

-- name: CreateCliente :one
INSERT INTO clientes(email, password) VALUES ($1, $2) RETURNING *;

-- name: CreateArticulo :one
INSERT INTO articulos(titulo, descripcion, imagen) VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteCliente :exec
DELETE FROM clientes WHERE id = $1;

-- name: DeleteArticulo :exec
DELETE FROM articulos WHERE id = $1;

-- name: UpdateCliente :exec
UPDATE clientes SET email = $2, password = $3 WHERE id = $1;

-- name: UpdateArticulo :exec
UPDATE articulos SET titulo = $2, descripcion = $3, imagen = $4 WHERE id = $1;
