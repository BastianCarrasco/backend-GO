// main.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	// "os" // Ya no se necesita os si no leemos variables de entorno
	"time"

	"github.com/BastianCarrasco/backend-go/db"
	"github.com/BastianCarrasco/backend-go/repository"
	"github.com/BastianCarrasco/backend-go/usecase"
	"github.com/gorilla/mux"
	// "github.com/joho/godotenv" // Ya no se necesita godotenv
)

// === Todas las configuraciones como constantes aquí ===
const (
	// ¡ADVERTENCIA DE SEGURIDAD EXTREMA!
	// NUNCA, NUNCA, NUNCA guardes credenciales sensibles como esta directamente en tu código fuente
	// en un entorno de producción o si tu código será público (ej. GitHub).
	// ESTA CONFIGURACIÓN ES SÓLO PARA PROPÓSITOS EDUCATIVOS O DEMOSTRATIVOS,
	// NO PARA PRODUCCIÓN. Para proyectos reales, USA SIEMPRE VARIABLES DE ENTORNO
	// o un sistema de gestión de secretos.

	appMongoURI       = "mongodb://mongo:eiyVLCEJjPpFKhkxPCkmXDVSFEqhrwGS@switchback.proxy.rlwy.net:52692"
	appDatabaseName   = "test"        // Nombre de tu base de datos (según tu ejemplo)
	appCollectionName = "Pruebas"     // Nombre de tu colección (según tu ejemplo)
	appPort           = ":3000"       // Puerto fijo directamente en el código
)

// La función init() ya no es necesaria si no se cargan variables de entorno
// ni se inicializan variables globales que dependen de ellas.
// Si aún tienes lógica en init() para otras cosas, déjala.
// De lo contrario, puedes eliminarla por completo.
func init() {
    // Si tu init() ya no hace nada, puedes eliminar este bloque.
    // Lo dejo comentado solo para fines de claridad sobre lo que se removió.
	// if err := godotenv.Load(); err != nil {
	// 	log.Println("No se encontró el archivo .env, intentando leer variables de entorno del sistema.")
	// }
	//
	// mongoURI = os.Getenv("MONGO_URI")
	// if mongoURI == "" {
	// 	log.Fatal("La variable de entorno MONGO_URI no está configurada. Por favor, establécela en .env o como variable del sistema.")
	// }
	//
	// databaseName = os.Getenv("DB_NAME")
	// if databaseName == "" {
	// 	log.Fatal("La variable de entorno DB_NAME no está configurada.")
	// }
	//
	// collectionName = os.Getenv("COLLECTION_NAME")
	// if collectionName == "" {
	// 	log.Fatal("La variable de entorno COLLECTION_NAME no está configurada.")
	// }
}

func main() {
	// 1. Conexión a MongoDB
	// Usamos appMongoURI directamente
	client, err := db.ConnectDB(appMongoURI)
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
	// Usamos appDatabaseName y appCollectionName directamente
	razaCollection := client.Database(appDatabaseName).Collection(appCollectionName)

	// 2. Inyección de Dependencias
	razaRepo := &repository.RazaRepository{MongoCollection: razaCollection}
	razaUseCase := usecase.NewRazaUseCase(*razaRepo)

	// 3. Configuración del Router (Gorilla Mux)
	router := mux.NewRouter()

	// Endpoints para Raza
	router.HandleFunc("/razas", getAllRacesHandler(razaUseCase)).Methods("GET")
	router.HandleFunc("/razas/{id}", getRazaByIDHandler(razaUseCase)).Methods("GET")

	// 4. Iniciar el Servidor HTTP
	log.Printf("Servidor escuchando en el puerto %s", appPort) // Usa la constante appPort
	log.Fatal(http.ListenAndServe(appPort, router))           // Usa la constante appPort
}

// ... los handlers (getAllRacesHandler, getRazaByIDHandler) permanecen iguales ...
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