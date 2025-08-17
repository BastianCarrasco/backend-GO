package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BastianCarrasco/backend-GO/controllers"
	"github.com/BastianCarrasco/backend-GO/db"
	"github.com/joho/godotenv"
)

// homeHandler responde a la ruta "/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	link := os.Getenv("LINK")
	log.Println("Acceso a / - LINK:", link)
	fmt.Fprintf(w, "¡Servidor web Go funcionando! Conectado a DB.")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se pudo cargar el archivo .env.")
	}

	db.ConnectDB()
	defer db.DisconnectDB()

	// Registrar manejadores de rutas
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/proyectos", controllers.GetProyectosHandler)
	// Para obtener por ID, la ruta debe terminar en "/" si el ID es un segmento.
	// Esto significa que GetProyectoByIDHandler manejará cualquier ruta que empiece con /proyectos/
	http.HandleFunc("/proyectos/", controllers.GetProyectoByIDHandler) // <--- ¡Nueva ruta para GET por ID!

	port := ":" + os.Getenv("PORT")
	if os.Getenv("PORT") == "" {
		port = ":8080"
	}

	fmt.Println("Servidor Go conectado correctamente. Escuchando en el puerto", port)
	log.Fatal(http.ListenAndServe(port, nil))
}