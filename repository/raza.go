package repository

import (
	"context" // Necesario para context.TODO()

	"github.com/BastianCarrasco/backend-go/model"
	"go.mongodb.org/mongo-driver/bson" // Necesario para bson.M y bson.D
	"go.mongodb.org/mongo-driver/mongo"
)

type RazaRepository struct {
	MongoCollection *mongo.Collection
}

func (r *RazaRepository) FindRazaByID(id string) (*model.Raza, error) { // Renombrado a PascalCase para ser exportable
	var raza model.Raza
	err := r.MongoCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&raza)
	if err != nil {
		// Manejo específico para cuando no se encuentra el documento
		if err == mongo.ErrNoDocuments {
			return nil, nil // Retorna nil, nil si no se encuentra la raza, indicando que no hubo error pero tampoco resultados.
		}
		return nil, err
	}
	return &raza, nil
}

func (r *RazaRepository) GetAllRaces() ([]model.Raza, error) {
	results, err := r.MongoCollection.Find(context.TODO(), bson.D{}) // bson.D{} para todos los documentos
	if err != nil {
		return nil, err
	}

	var razas []model.Raza
	// Itera sobre los resultados y decodifica cada documento en un objeto Raza
	for results.Next(context.TODO()) {
		var raza model.Raza
		err := results.Decode(&raza)
		if err != nil {
			return nil, err // Si hay un error al decodificar un documento, retorna el error
		}
		razas = append(razas, raza)
	}

	// Verifica si hubo algún error durante la iteración (después de Next())
	if err := results.Err(); err != nil {
		return nil, err
	}

	return razas, nil
}