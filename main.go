package main

import (
	"encoding/json" // Para trabajar con JSON
	"fmt"           // Para formatear cadenas de texto e imprimir
	"log"           // Para logging de errores
	"net/http"      // Para el servidor HTTP
)

// Definición de una estructura para simular los datos de la empresa
// Las etiquetas `json:"nombreCampo"` son importantes para la serialización/deserialización a JSON
type Empresa struct {
	ID                 string   `json:"_id"`
	Nombre             string   `json:"nombre"`
	Apellido           string   `json:"apellido"`
	EmpresaOrganizacion string   `json:"empresaOrganizacion"`
	AreaTrabajo        string   `json:"areaTrabajo"`
	CorreoElectronico  string   `json:"correoElectronico"`
	NumeroTelefono     string   `json:"numeroTelefono"`
	ContactoWeb        string   `json:"contactoWeb"`
	VinculoPUCV        []string `json:"vinculoPUCV"`
	ActividadesServicios string   `json:"actividadesServicios"`
	Desafio1           string   `json:"desafio1"`
	Desafio2           string   `json:"desafio2"`
	Desafio3           string   `json:"desafio3"`
	InteresInformacion string   `json:"interesInformacion"`
	CreatedAt          string   `json:"createdAt"`
	UpdatedAt          string   `json:"updatedAt"`
}

// Estructura para la respuesta de la API (simulando el formato de tu frontend)
type ApiResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Data    []Empresa `json:"data"`
}

// Datos de empresas de ejemplo
var empresasData = []Empresa{
	{
		ID:                  "60d5ec49f8c6d40015b6d5e7",
		Nombre:              "Juan",
		Apellido:            "Pérez",
		EmpresaOrganizacion: "Tech Solutions S.A.",
		AreaTrabajo:         "Desarrollo de Software",
		CorreoElectronico:   "juan.perez@techsol.com",
		NumeroTelefono:      "+56912345678",
		ContactoWeb:         "http://www.techsol.com",
		VinculoPUCV:         []string{"Ex-alumno Ingeniería Informática", "Mentor Startup PUCV"},
		ActividadesServicios: "Desarrollo de aplicaciones web y móviles.",
		Desafio1:            "Escalar nuestra infraestructura en la nube.",
		Desafio2:            "Integrar IA en nuestros productos existentes.",
		Desafio3:            "Atraer talento especializado en ciberseguridad.",
		InteresInformacion:  "si",
		CreatedAt:           "2023-01-15T10:00:00Z",
		UpdatedAt:           "2023-01-15T10:00:00Z",
	},
	{
		ID:                  "60d5ec49f8c6d40015b6d5e8",
		Nombre:              "María",
		Apellido:            "Gómez",
		EmpresaOrganizacion: "Innovate Consultores",
		AreaTrabajo:         "Consultoría de Negocios",
		CorreoElectronico:   "maria.gomez@innovate.cl",
		NumeroTelefono:      "+56987654321",
		ContactoWeb:         "http://www.innovate.cl",
		VinculoPUCV:         []string{"Docente invitada Finanzas"},
		ActividadesServicios: "Asesoría estratégica y gestión de proyectos.",
		Desafio1:            "Expansión a mercados internacionales.",
		Desafio2:            "Optimización de procesos internos con tecnología.",
		Desafio3:            "", // Puede estar vacío
		InteresInformacion:  "no",
		CreatedAt:           "2023-02-01T11:30:00Z",
		UpdatedAt:           "2023-02-01T11:30:00Z",
	},
	{
		ID:                  "60d5ec49f8c6d40015b6d5e9",
		Nombre:              "Carlos",
		Apellido:            "Rodríguez",
		EmpresaOrganizacion: "Ecodevelop Sostenible",
		AreaTrabajo:         "Energías Renovables",
		CorreoElectronico:   "carlos.rodriguez@ecodevelop.cl",
		NumeroTelefono:      "+56923456789",
		ContactoWeb:         "http://www.ecodevelop.cl",
		VinculoPUCV:         []string{"Colaborador de Investigación"},
		ActividadesServicios: "Implementación de soluciones solares y eólicas.",
		Desafio1:            "Obtener financiamiento para proyectos de gran escala.",
		Desafio2:            "Desarrollar nuevas tecnologías de almacenamiento de energía.",
		Desafio3:            "",
		InteresInformacion:  "si",
		CreatedAt:           "2023-03-10T09:15:00Z",
		UpdatedAt:           "2023-03-10T09:15:00Z",
	},
}

// homeHandler responde a la ruta "/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "¡Bienvenido a tu backend simple en Go!")
}

// saludoHandler responde a la ruta "/api/saludo"
func saludoHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el parámetro 'nombre' de la query string
	nombre := r.URL.Query().Get("nombre")

	if nombre == "" {
		nombre = "mundo" // Valor por defecto si no se proporciona 'nombre'
	}

	fmt.Fprintf(w, "¡Hola, %s! Desde Go.", nombre)
}

// empresasHandler responde a la ruta "/api/empresas"
func empresasHandler(w http.ResponseWriter, r *http.Request) {
	// Configurar el encabezado para indicar que la respuesta es JSON
	w.Header().Set("Content-Type", "application/json")
	// Configurar CORS para permitir peticiones desde cualquier origen (para desarrollo)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Si es una petición OPTIONS (preflight CORS), solo enviar los headers y salir
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Crear la respuesta con el formato ApiResponse
	response := ApiResponse{
		Status:  "success",
		Message: "Empresas obtenidas correctamente",
		Data:    empresasData,
	}

	// Codificar la estructura a JSON y enviarla como respuesta
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	// Registrar los manejadores de ruta
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/saludo", saludoHandler)
	http.HandleFunc("/api/empresas", empresasHandler) // Nuevo endpoint para las empresas

	// Definir el puerto en el que el servidor escuchará
	port := ":8080"
	fmt.Printf("Servidor Go escuchando en http://localhost%s\n", port)

	// Iniciar el servidor HTTP
	// log.Fatal detendrá el programa si hay un error al iniciar el servidor
	log.Fatal(http.ListenAndServe(port, nil))
}