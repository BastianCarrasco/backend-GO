package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	// Necesario para convertir int a string
	// Necesario para strings.TrimPrefix
	_ "github.com/BastianCarrasco/backend-GO/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/BastianCarrasco/backend-GO/controllers"
	"github.com/BastianCarrasco/backend-GO/db"
	"github.com/joho/godotenv"
)

// homeHandler responde a la ruta "/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	link := os.Getenv("LINK") // Esto es la URI de la DB
	log.Println("Acceso a / - LINK (DB URI):", link)
	fmt.Fprintf(w, "¡Servidor web Go funcionando! Conectado a DB.")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se pudo cargar el archivo .env.")
	}

	db.ConnectDB()
	defer db.DisconnectDB()

	// Definir el puerto que el servidor escuchará (incluye el ":")
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080" // Puerto por defecto si no se especifica
	}
	listenAddr := ":" + portStr // Ejemplo: ":8080"

	// Construir la URL base de la API y de Swagger UI
	// Esto asume que la API se expone en "localhost" en desarrollo.
	// En un entorno desplegado, esto debería ser el dominio público.
	apiHost := "localhost:" + portStr // Ejemplo: "localhost:8080"
	swaggerDocURL := fmt.Sprintf("http://%s/swagger/doc.json", apiHost)

	// Registrar manejadores de rutas
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/proyectos", controllers.GetProyectosHandler)
	http.HandleFunc("/proyectos/", controllers.GetProyectoByIDHandler)

	// --- Configuración para Swagger UI ---
	// @host y @schemes en las anotaciones de Swagger (arriba) deben coincidir con `apiHost`
	// y el protocolo ("http" o "https").
	// Regenera la documentación con 'swag init' después de cualquier cambio en `@host` o `@schemes`.
	http.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL(swaggerDocURL), // La URL de donde se cargará el JSON/YAML de la API
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))
	// ------------------------------------

	fmt.Println("Servidor Go conectado correctamente. Escuchando en el puerto", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

// @title API de Proyectos UNAB
// @version 1.0
// @description Esta es la API para la gestión de proyectos de investigación de la UNAB.
// @termsOfService http://swagger.io/terms/

// @contact.name Equipo de Desarrollo
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080 
// @BasePath /
// @schemes http // Si usas HTTPS en producción, añade "https" aquí también
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization