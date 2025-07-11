package main

import (
	"fmt"      // Para formatear cadenas de texto e imprimir
	"log"      // Para logging de errores
	"net/http" // Para el servidor HTTP
)

// homeHandler responde a la ruta "/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "¡Servidor web Go funcionando!")
}

func main() {
	// Registrar el manejador de ruta para la raíz
	http.HandleFunc("/", homeHandler)

	// Definir el puerto en el que el servidor escuchará
	port := ":8080"

	// Mensaje simplificado para la consola
	fmt.Println("Servidor Go conectado correctamente. Escuchando en el puerto", port)

	// Iniciar el servidor HTTP
	// log.Fatal detendrá el programa si hay un error al iniciar el servidor
	log.Fatal(http.ListenAndServe(port, nil))
}