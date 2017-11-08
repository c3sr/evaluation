package evaluation

import (
	"time"

	"github.com/rai-project/database"
	"github.com/rai-project/database/mongodb"
	"gopkg.in/mgo.v2/bson"
)

type ModelAccuracy struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt time.Time     `json:"created_at"  bson:"created_at"`
	Top1      float64
	Top5      float64
}

func (ModelAccuracy) TableName() string {
	return "model_accuracy"
}

type ModelAccuracyCollection struct {
	*mongodb.MongoTable
}

func NewModelAccuracyCollection(db database.Database) (*ModelAccuracyCollection, error) {
	tbl, err := mongodb.NewTable(db, ModelAccuracy{}.TableName())
	if err != nil {
		return nil, err
	}
	tbl.Create(nil)

	return &ModelAccuracyCollection{
		MongoTable: tbl.(*mongodb.MongoTable),
	}, nil
}

func (m *ModelAccuracyCollection) Close() error {
	return nil
}