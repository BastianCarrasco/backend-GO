package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref" // Opcional, para verificar la conexión
)

// Client es la instancia global del cliente de MongoDB
var Client *mongo.Client

// ConnectDB establece la conexión con MongoDB
func ConnectDB() {
	// Obtener la URI de MongoDB desde las variables de entorno
	mongoURI := os.Getenv("LINK")
	if mongoURI == "" {
		log.Fatal("Error: La variable de entorno 'LINK' (MONGO_URI) no está configurada.")
	}

	// Establecer el contexto con un tiempo límite para la conexión
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Asegura que el contexto se cancele al finalizar

	// Opciones del cliente de MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Conectarse a MongoDB
	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error al conectar a MongoDB: %v", err)
	}

	// Opcional: Hacer un ping para verificar la conexión
	err = Client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Error al hacer ping a MongoDB: %v", err)
	}

	log.Println("MongoDB conectado exitosamente a:", mongoURI)
}

// GetCollection obtiene una colección específica de la base de datos
func GetCollection(collectionName string) *mongo.Collection {
	if Client == nil {
		log.Fatal("Error: El cliente de MongoDB no está inicializado. Llama a ConnectDB() primero.")
	}
	// Reemplaza "Cartera_Mongo" con el nombre real de tu base de datos
	return Client.Database("Cartera_Mongo").Collection(collectionName)
}

// DisconnectDB cierra la conexión con MongoDB
func DisconnectDB() {
	if Client == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Client.Disconnect(ctx)
	if err != nil {
		log.Printf("Error al desconectar de MongoDB: %v", err)
	} else {
		log.Println("MongoDB desconectado.")
	}
}