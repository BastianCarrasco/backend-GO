// main.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os" // Necesario para os.Getenv
	"time"

	"github.com/BastianCarrasco/backend-go/db"
	"github.com/BastianCarrasco/backend-go/repository"
	"github.com/BastianCarrasco/backend-go/usecase"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv" // Importar godotenv
)

// Estas constantes ya no necesitan valores predeterminados.
// Se obtendrán de las variables de entorno.
var (
	mongoURI       string
	databaseName   string
	collectionName string
	port           string
)

func init() {
	// init() se ejecuta automáticamente antes de main()
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró el archivo .env, intentando leer variables de entorno del sistema.")
		// Si no se encuentra .env, el programa intentará leer directamente de las variables de entorno del sistema.
		// Esto es útil para entornos de despliegue donde .env no existe y las variables están preconfiguradas.
	}

	// Obtener los valores de las variables de entorno
	mongoURI = os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("La variable de entorno MONGO_URI no está configurada. Por favor, establécela en .env o como variable del sistema.")
	}

	databaseName = os.Getenv("DB_NAME")
	if databaseName == "" {
		log.Fatal("La variable de entorno DB_NAME no está configurada.")
	}

	collectionName = os.Getenv("COLLECTION_NAME")
	if collectionName == "" {
		log.Fatal("La variable de entorno COLLECTION_NAME no está configurada.")
	}

	port = os.Getenv("PORT")
	if port == "" {
		log.Println("La variable de entorno PORT no está configurada, usando puerto por defecto: :3000")
		port = ":8080" // Puerto por defecto si no se especifica en .env
	}
}

func main() {
	// 1. Conexión a MongoDB usando la función del paquete 'db'
	// db.ConnectDB() necesita la URI, por eso la pasamos aquí
	client, err := db.ConnectDB(mongoURI) // Pasar mongoURI a ConnectDB
	if err != nil {
		log.Fatalf("Fallo crítico al conectar a la base de datos: %v", err)
	}
	defer func() {
		disconnectCtx, cancelDisconnect := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelDisconnect()
		if err = client.Disconnect(disconnectCtx); err != nil {
			log.Fatalf("Error al desconectar de MongoDB: %v", err)
		}
		log.Println("Desconectado de MongoDB.")
	}()

	// Obtener la colección
	razaCollection := client.Database(databaseName).Collection(collectionName)

	// 2. Inyección de Dependencias
	razaRepo := &repository.RazaRepository{MongoCollection: razaCollection}
	razaUseCase := usecase.NewRazaUseCase(*razaRepo)

	// 3. Configuración del Router (Gorilla Mux)
	router := mux.NewRouter()

	// Endpoints para Raza
	router.HandleFunc("/razas", getAllRacesHandler(razaUseCase)).Methods("GET")
	router.HandleFunc("/razas/{id}", getRazaByIDHandler(razaUseCase)).Methods("GET")

	// 4. Iniciar el Servidor HTTP
	log.Printf("Servidor escuchando en el puerto %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}

// ... los handlers (getAllRacesHandler, getRazaByIDHandler) permanecen iguales ...
// Handler para obtener todas las razas
func getAllRacesHandler(uc usecase.RazaUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		razas, err := uc.GetAllRaces(ctx)
		if err != nil {
			log.Printf("Error al obtener todas las razas: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(razas); err != nil {
			log.Printf("Error al codificar respuesta JSON: %v", err)
			http.Error(w, "Error al serializar respuesta", http.StatusInternalServerError)
		}
	}
}

// Handler para obtener una raza por ID
func getRazaByIDHandler(uc usecase.RazaUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		raza, err := uc.GetRazaByID(ctx, id)
		if err != nil {
			log.Printf("Error al obtener raza por ID %s: %v", id, err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		if raza == nil {
			http.Error(w, "Raza no encontrada", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(raza); err != nil {
			log.Printf("Error al codificar respuesta JSON: %v", err)
			http.Error(w, "Error al serializar respuesta", http.StatusInternalServerError)
		}
	}
}