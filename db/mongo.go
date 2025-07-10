// db/mongo.go
package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ConnectDB establece la conexión a MongoDB y la devuelve.
// Ahora toma mongoURI como argumento.
func ConnectDB(mongoURI string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Printf("Error al conectar a MongoDB: %v", err)
		return nil, err
	}

	// Ping para verificar la conexión
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		client.Disconnect(ctx) // Asegurarse de desconectar si el ping falla
		log.Printf("Error al hacer ping a MongoDB: %v", err)
		return nil, err
	}

	log.Println("Conectado exitosamente a MongoDB!")
	return client, nil
}