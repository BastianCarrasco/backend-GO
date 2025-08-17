package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client

func ConnectDB() {
	mongoURI := os.Getenv("LINK")
	if mongoURI == "" {
		log.Fatal("Error: La variable de entorno 'LINK' (MONGO_URI) no está configurada.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)

	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error al conectar a MongoDB: %v", err)
	}

	err = Client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Error al hacer ping a MongoDB: %v", err)
	}

	log.Println("MongoDB conectado exitosamente a:", mongoURI)
}

// GetCollection obtiene una colección específica de la base de datos de forma genérica
func GetCollection(collectionName string) *mongo.Collection {
	if Client == nil {
		log.Fatal("Error: El cliente de MongoDB no está inicializado. Llama a ConnectDB() primero.")
	}
	return Client.Database("CARTERA").Collection(collectionName) 
}

// GetProyectosCollection obtiene específicamente la colección "PROYECTOS"
// (Es redundante si siempre vas a usar GetCollection, pero puede ser útil para claridad en el código)
func GetProyectosCollection() *mongo.Collection {
	return GetCollection("PROYECTOS") // Reutiliza la función GetCollection
}

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