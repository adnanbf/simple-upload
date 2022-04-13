package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadData struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Name         string             `json:"-" bson:"name"`
	OriginalName string             `json:"original_name" bson:"original_name"`
	FileSize     int64              `json:"file_size" bson:"file_size"`
	Ext          string             `json:"ext" bson:"ext"`
	Mime         string             `json:"mime" bson:"mime"`
	CreatedAt    *time.Time         `json:"created_at" bson:"created_at"`
}

func NewUploadData() *UploadData {
	tNow := time.Now()
	return &UploadData{
		ID:        primitive.NewObjectID(),
		CreatedAt: &tNow,
	}
}

func (r *UploadData) Table() string {
	return "upload"
}
