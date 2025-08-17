// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os" // Asegúrate de importar os

	"github.com/joho/godotenv"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    link := os.Getenv("LINK") // Ahora LINK ya estará cargado
    log.Println("LINK dentro de homeHandler:", link)
    fmt.Fprintf(w, "¡Servidor web Go funcionando! LINK: %s", link)
}

func main() {
    // Cargar variables de entorno una única vez al inicio de main
    err := godotenv.Load()
    if err != nil {
        log.Println("Advertencia: No se pudo cargar el archivo .env")
        // No es Fatal aquí si no es estrictamente necesario,
        // ya que podría haber vars de entorno ya configuradas en el sistema.
    }

    link := os.Getenv("LINK")
    log.Println("LINK al iniciar el servidor:", link)

    http.HandleFunc("/", homeHandler)

    port := ":8080"
    fmt.Println("Servidor Go conectado correctamente. Escuchando en el puerto", port)
    log.Fatal(http.ListenAndServe(port, nil))
}