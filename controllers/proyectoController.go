package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings" // Necesario para strings.Split
	"time"

	"github.com/BastianCarrasco/backend-GO/db"
	"github.com/BastianCarrasco/backend-GO/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo" // Asegúrate de que mongo esté importado
)

// setCORSHeaders es una función auxiliar para establecer los encabezados CORS
// Esto ayuda a mantener el código DRY (Don't Repeat Yourself)
func setCORSHeaders(w http.ResponseWriter, methods string) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Permite cualquier origen (para desarrollo)
	w.Header().Set("Access-Control-Allow-Methods", methods)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// GetProyectosHandler obtiene todos los proyectos de la colección
func GetProyectosHandler(w http.ResponseWriter, r *http.Request) {
	// Manejar preflight OPTIONS request para CORS
	if r.Method == http.MethodOptions {
		setCORSHeaders(w, "GET, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	setCORSHeaders(w, "GET") // Establece los encabezados CORS para GET

	collection := db.GetCollection("PROYECTOS") // O db.GetProyectosCollection()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error al buscar proyectos: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var proyectos []models.Proyecto
	if err = cursor.All(ctx, &proyectos); err != nil {
		log.Printf("Error al decodificar proyectos: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proyectos)
}

// GetProyectoByIDHandler obtiene un proyecto específico por su ID
func GetProyectoByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Manejar preflight OPTIONS request para CORS
	if r.Method == http.MethodOptions {
		setCORSHeaders(w, "GET, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	setCORSHeaders(w, "GET") // Establece los encabezados CORS para GET

	// Extraer el ID de la URL
	// Usamos strings.Split para ser explícitos.
	// r.URL.Path será algo como "/proyectos/689d6c..."
	pathSegments := strings.Split(r.URL.Path, "/")
	// pathSegments[0] será "", pathSegments[1] será "proyectos", pathSegments[2] será el ID
	if len(pathSegments) < 3 || pathSegments[2] == "" { // Verifica que haya un ID en la posición 2
		http.Error(w, "ID de proyecto no proporcionado en la URL", http.StatusBadRequest)
		return
	}
	proyectoIDStr := pathSegments[2]

	// Convertir el ID de string a primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(proyectoIDStr)
	if err != nil {
		http.Error(w, "ID de proyecto inválido", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection("PROYECTOS") // O db.GetProyectosCollection()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var proyecto models.Proyecto
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&proyecto)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Proyecto no encontrado", http.StatusNotFound)
			return
		}
		log.Printf("Error al buscar proyecto por ID: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proyecto)
}

// NOTA: Para POST, PUT, DELETE, añadirías funciones similares en este archivo.