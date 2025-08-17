package main

import (
	"fmt"
	"log"
	"net/http"
	"os" // Necesario para os.Getenv si no lo tienes

	"github.com/BastianCarrasco/backend-GO/db" // Importa tu paquete db

	"github.com/joho/godotenv"
)

// homeHandler responde a la ruta "/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Ya no es necesario cargar .env aquí, se hace una vez en main()
	link := os.Getenv("LINK") // 'LINK' ya estará cargado
	log.Println("Acceso a / - LINK:", link)

	// Opcional: Puedes verificar el estado de la conexión a la DB aquí
	// if db.Client == nil {
	//     fmt.Fprintf(w, "¡Servidor web Go funcionando! Pero la DB no está conectada.")
	//     return
	// }
	// err := db.Client.Ping(context.Background(), nil)
	// if err != nil {
	//     fmt.Fprintf(w, "¡Servidor web Go funcionando! DB: Error de ping - %v", err)
	//     return
	// }

	fmt.Fprintf(w, "¡Servidor web Go funcionando! Conectado a DB.")
}

func main() {
	// 1. Cargar variables de entorno una única vez al inicio de main
	err := godotenv.Load()
	if err != nil {
		// Loguea una advertencia si .env no se encuentra, pero no es fatal si no es estrictamente necesario
		log.Println("Advertencia: No se pudo cargar el archivo .env. Asegúrate de que las variables de entorno estén configuradas.")
	}

	// 2. Conectar a MongoDB
	db.ConnectDB()
	// Asegúrate de desconectar la DB cuando la aplicación se cierre (ej. Ctrl+C)
	defer db.DisconnectDB() // Esto se ejecutará justo antes de que main termine

	// Registrar el manejador de ruta para la raíz
	http.HandleFunc("/", homeHandler)

	// Definir el puerto en el que el servidor escuchará
	port := ":" + os.Getenv("PORT") // Usa una variable de entorno para el puerto si la tienes
	if os.Getenv("PORT") == "" {
		port = ":8080" // Puerto por defecto si no se define en .env
	}


	// Mensaje simplificado para la consola
	fmt.Println("Servidor Go conectado correctamente. Escuchando en el puerto", port)

	// Iniciar el servidor HTTP
	// log.Fatal detendrá el programa si hay un error al iniciar el servidor
	log.Fatal(http.ListenAndServe(port, nil))
}