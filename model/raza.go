// model/raza.go
package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Unidad representa la estructura de unidad_mas_simple
type Unidad struct {
	Nombre string `json:"nombre" bson:"nombre"`
	Funcion string `json:"funcion" bson:"funcion"`
}

// Raza representa la estructura de una raza de mascota o facción en la base de datos.
type Raza struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Nombre         string             `json:"nombre" bson:"nombre"`
	Caracteristicas []string           `json:"caracteristicas" bson:"caracteristicas"` // Un array de strings
	UnidadMasSimple Unidad             `json:"unidad_mas_simple" bson:"unidad_mas_simple"` // Un objeto anidado
}