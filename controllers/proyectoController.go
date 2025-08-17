package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/BastianCarrasco/backend-GO/db"
	"github.com/BastianCarrasco/backend-GO/models" // Asegúrate de que models.Proyecto esté bien definido aquí

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// setCORSHeaders es una función auxiliar para establecer los encabezados CORS
func setCORSHeaders(w http.ResponseWriter, methods string) {
	w.Header().Set("Access-Control-Allow-Origin", "*") 
	w.Header().Set("Access-Control-Allow-Methods", methods)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// @Summary Obtener todos los proyectos
// @Description Obtiene una lista de todos los proyectos de investigación disponibles.
// @Tags Proyectos
// @Accept json
// @Produce json
// @Success 200 {array} models.Proyecto "Lista de proyectos obtenida exitosamente"
// @Failure 500 {string} string "Error interno del servidor al buscar proyectos"
// @Router /proyectos [get]
func GetProyectosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		setCORSHeaders(w, "GET, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	setCORSHeaders(w, "GET")

	collection := db.GetCollection("PROYECTOS") 

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

// @Summary Obtener un proyecto por ID
// @Description Obtiene los detalles de un proyecto específico utilizando su ID.
// @Tags Proyectos
// @Accept json
// @Produce json
// @Param id path string true "ID del Proyecto (MongoDB ObjectID)"
// @Success 200 {object} models.Proyecto "Detalles del proyecto obtenidos exitosamente"
// @Failure 400 {string} string "ID de proyecto inválido"
// @Failure 404 {string} string "Proyecto no encontrado"
// @Failure 500 {string} string "Error interno del servidor al buscar proyecto"
// @Router /proyectos/{id} [get]
func GetProyectoByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		setCORSHeaders(w, "GET, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	setCORSHeaders(w, "GET")

	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) < 3 || pathSegments[2] == "" {
		http.Error(w, "ID de proyecto no proporcionado en la URL", http.StatusBadRequest)
		return
	}
	proyectoIDStr := pathSegments[2]

	objID, err := primitive.ObjectIDFromHex(proyectoIDStr)
	if err != nil {
		http.Error(w, "ID de proyecto inválido", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection("PROYECTOS") 

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