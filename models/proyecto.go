package models

import (
	"time" // Para el tipo fecha_postulacion

	"go.mongodb.org/mongo-driver/bson/primitive" // Para el tipo _id de MongoDB
)

// Academicos representa la estructura de cada académico dentro del array "academicos"
type Academico struct {
	Nombre   string `bson:"nombre"`
	APaterno string `bson:"a_paterno"`
	AMaterno string `bson:"a_materno"`
}

// Proyecto representa la estructura de un documento en la colección PROYECTOS
type Proyecto struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"` // _id de MongoDB, omitempty para inserciones
	Nombre             string             `bson:"nombre"`
	Academicos         []Academico        `bson:"academicos"`       // Un slice de Academicos
	Estudiantes        []interface{}      `bson:"estudiantes"`      // Parece ser un array vacío, puedes usar interface{} o un tipo más específico si cambia
	Monto              int                `bson:"monto"`
	FechaPostulacion   time.Time          `bson:"fecha_postulacion"` // Usamos time.Time para fechas ISO
	Unidad             string             `bson:"unidad"`
	Tematica           string             `bson:"tematica"`
	Estatus            string             `bson:"estatus"`
	Convocatoria       string             `bson:"convocatoria"`
	TipoConvocatoria   string             `bson:"tipo_convocatoria"`
	InstConv           string             `bson:"inst_conv"`
	DetalleApoyo       string             `bson:"detalle_apoyo"`
	Apoyo              string             `bson:"apoyo"`
	IDKTH              *string            `bson:"id_kth,omitempty"`  // Usamos *string para campos que pueden ser null
	Comentarios        string             `bson:"comentarios"`
}