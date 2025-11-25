package main

import (
	"fmt"
	"log"
	"net/http"

	"catalogo.com/proyecto-catalogo/handlers"

	bd "catalogo.com/proyecto-catalogo/bd"
	sqlc "catalogo.com/proyecto-catalogo/bd/sqlc"
	_ "github.com/lib/pq"
)

func main() {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Conectar la base de datos
	conn, err := bd.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Crear queries de sqlc
	queries := sqlc.New(conn)

	h := handlers.NewHandler(queries)

	http.HandleFunc("/", h.TemplHandler)

	http.HandleFunc("/clientes", h.CreateClienteForm)

	http.HandleFunc("/articulos", h.CreateArticuloForm)

	http.HandleFunc("/clientes/", h.HandlerDeleteCliente)

	http.HandleFunc("/articulos/{id}", h.HandlerDeleteArticulo)

	port1 := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port1)

	// Start server
	err = http.ListenAndServe(port1, nil)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %s\n", err)
	}

}
