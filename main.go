package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/BastianCarrasco/backend-GO/docs"
	_ "github.com/BastianCarrasco/backend-GO/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/BastianCarrasco/backend-GO/controllers"
	"github.com/BastianCarrasco/backend-GO/db"
	"github.com/joho/godotenv"
)

// homeHandler responde a la ruta "/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	link := os.Getenv("LINK")
	log.Println("Acceso a / - LINK (DB URI):", link)
	fmt.Fprintf(w, "¡Servidor web Go funcionando! Conectado a DB.")
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

// @host localhost:8080 // ESTA ANOTACIÓN SERÁ SOBREESCRITA POR main()
// @BasePath /
// @schemes http // ESTA ANOTACIÓN SERÁ SOBREESCRITA POR main()
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// 1. Cargar variables de entorno
	// godotenv.Load() solo carga variables si no existen ya en el entorno.
	// Esto es ideal, ya que Railway inyectará sus propias variables.
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se pudo cargar el archivo .env.")
	}

	// 2. Conectar a MongoDB
	db.ConnectDB()
	defer db.DisconnectDB() // Asegura que la conexión a la DB se cierre al finalizar

	// 3. Determinar el puerto de escucha
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080" // Puerto por defecto para desarrollo local si no se define
	}
	listenAddr := ":" + portStr // Ejemplo: ":8080" o ":19432"

	// 4. Determinar el host para Swagger (para desarrollo local o despliegue en Railway)
	var apiHost string
	// Si RAILWAY_PUBLIC_URL existe, significa que estamos en Railway
	if railwayURL := os.Getenv("RAILWAY_PUBLIC_URL"); railwayURL != "" {
		apiHost = strings.TrimPrefix(railwayURL, "https://") // Quita el "https://"
		apiHost = strings.TrimSuffix(apiHost, "/") // Quita barras finales
		// Nota: En Railway, la URL pública no incluye el puerto, ya que va por 80/443
		// Si RAILWAY_PUBLIC_URL incluye el puerto, lo dejamos, pero no suele ser el caso.
	} else {
		// Para desarrollo local
		apiHost = "localhost:" + portStr
	}

	// 5. Construir la URL completa para el swagger.json
	// Usamos HTTPS si estamos en Railway (o si la URL pública lo indica), de lo contrario HTTP
	swaggerProtocol := "http"
	if strings.HasPrefix(os.Getenv("RAILWAY_PUBLIC_URL"), "https://") {
		swaggerProtocol = "https"
	}
	swaggerDocURL := fmt.Sprintf("%s://%s/swagger/doc.json", swaggerProtocol, apiHost)


	// 6. Configurar el "Host" y "Schemes" de Swagger programáticamente
	// Esto es crucial para que Swagger UI genere las URLs correctas para las peticiones.
	// La anotación @host arriba es un valor por defecto para 'swag init',
	// pero en runtime la sobrescribimos aquí.
	docs.SwaggerInfo.Host = apiHost
	docs.SwaggerInfo.Schemes = []string{swaggerProtocol}

	// 7. Registrar manejadores de rutas
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/proyectos", controllers.GetProyectosHandler)
	http.HandleFunc("/proyectos/", controllers.GetProyectoByIDHandler)

	// --- Configuración para Swagger UI ---
	http.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL(swaggerDocURL), // Usa la URL dinámica
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))
	// ------------------------------------

	fmt.Println("Servidor Go conectado correctamente. Escuchando en el puerto", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}