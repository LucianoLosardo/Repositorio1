package handlers

import (
	"context"
	"net/http"
	"strconv"

	"strings"

	bd "catalogo.com/proyecto-catalogo/bd/sqlc"
	"catalogo.com/proyecto-catalogo/views"
	"github.com/a-h/templ"
)

// structs
type Cliente struct { //se declara otro struct para evitar tere que recibir id (el del models lo pide)
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Articulo struct {
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	Imagen      string `json:"imagen"`
}

type Queries struct {
	queries *bd.Queries
}

func NewHandler(q *bd.Queries) *Queries {
	return &Queries{queries: q}
}

// handler principal
func (h *Queries) TemplHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}
	clients, err1 := h.queries.ListClientes(context.Background())
	articulos, err := h.queries.ListArticulos(context.Background())
	if err != nil && err1 != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	templ.Handler(views.IndexPage(clients, articulos)).ServeHTTP(w, r)
}

//////////////////////////////////
////////handlers articulo/////////
//////////////////////////////////

// validar que el articulo este bien formado
func validarArticulo(art Articulo) map[string]string {
	errors := make(map[string]string)

	if art.Titulo == "" {
		errors["titulo"] = "El titulo no puede ser vacio"
	}
	if art.Descripcion == "" {
		errors["descripcion"] = "La descripcion no puede ser vacia"
	}
	if art.Imagen == "" {
		errors["imagen"] = "La imagen no puede ser vacia"
	}
	return errors
}

// crear articulo desde formulario
func (h *Queries) CreateArticuloForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	// parsear los datos input del formulario para que esten disponibles
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario", http.StatusBadRequest)
		return
	}

	newArticulo := bd.CreateArticuloParams{
		Titulo:      r.FormValue("titulo"),
		Descripcion: r.FormValue("descripcion"),
		Imagen:      r.FormValue("imagen"),
	}

	// validacion (como estaba antes)
	validationErrors := validarArticulo(Articulo(newArticulo))
	if len(validationErrors) > 0 {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// insert
	_, err := h.queries.CreateArticulo(r.Context(), newArticulo)
	if err != nil {
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	articulos, err := h.queries.ListArticulos(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	views.ArticuloList(articulos).Render(r.Context(), w)
}

// /// HANDLER para /articulos/{id}  (DELETE)////
// funcion
func (h *Queries) deleteArticulo(w http.ResponseWriter, r *http.Request, id int32) {
	err := h.queries.DeleteArticulo(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//cargar la lista actualizada de articulos
	articulos, err := h.queries.ListArticulos(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//renderiza la lista actualizada
	views.ArticuloList(articulos).Render(r.Context(), w)

	w.WriteHeader(http.StatusOK)
}

// handler
func (h *Queries) HandlerDeleteArticulo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])

	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	h.deleteArticulo(w, r, int32(id))
}

//////////////////////////////////
//////// handlers cliente/////////
//////////////////////////////////

// validar que el cliente este bien formado
func validarCliente(client Cliente) map[string]string {
	errors := make(map[string]string)

	if client.Email == "" {
		errors["email"] = "El email es obligatorio"
	}
	if client.Password == "" {
		errors["password"] = "La contraseña es obligatoria"
	}

	return errors
}

// /// HANDLER para /clientes  (POST, insert)////
func (h *Queries) CreateClienteForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	// parsear los datos input del formulario para que esten disponibles
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario", http.StatusBadRequest)
		return
	}

	newCliente := bd.CreateClienteParams{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	// validacion (como estaba antes)
	validationErrors := validarCliente(Cliente(newCliente))
	if len(validationErrors) > 0 {
		http.Error(w, "Datos invalidos", http.StatusBadRequest)
		return
	}

	// insert
	_, err := h.queries.CreateCliente(r.Context(), newCliente)
	if err != nil {
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	clients, err := h.queries.ListClientes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	views.ClientList(clients).Render(r.Context(), w)
}

// /// HANDLER para /clientes/{id}  (DELETE)////
// funcion
func (h *Queries) deleteCliente(w http.ResponseWriter, r *http.Request, id int32) {
	err := h.queries.DeleteCliente(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//cargar la lista actualizada de clientes
	clients, err := h.queries.ListClientes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//rendderizar la lista actualizada
	views.ClientList(clients).Render(r.Context(), w)

	w.WriteHeader(http.StatusOK)
}

// handler
func (h *Queries) HandlerDeleteCliente(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])

	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	h.deleteCliente(w, r, int32(id))
}
