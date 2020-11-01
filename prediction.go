package evaluation

import (
	"time"

	"github.com/rai-project/database"
	"github.com/rai-project/database/mongodb"
	"github.com/c3sr/dlframework"
	"gopkg.in/mgo.v2/bson"
)

type InputPrediction struct {
	ID            bson.ObjectId        `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt     time.Time            `json:"created_at,omitempty"`
	InputID       string               `json:"input_id,omitempty"`
	InputIndex    int                  `json:"input_index,omitempty"`
	ExpectedLabel string               `json:"expected_label,omitempty"`
	Features      dlframework.Features `json:"features,omitempty"`
}

func (InputPrediction) TableName() string {
	return "input_prediction"
}

type InputPredictionCollection struct {
	*mongodb.MongoTable
}

func NewInputPredictionCollection(db database.Database) (*InputPredictionCollection, error) {
	tbl, err := mongodb.NewTable(db, InputPrediction{}.TableName())
	if err != nil {
		return nil, err
	}
	tbl.Create(nil)

	return &InputPredictionCollection{
		MongoTable: tbl.(*mongodb.MongoTable),
	}, nil
}

func (c *InputPredictionCollection) Find(as ...interface{}) ([]InputPrediction, error) {
	preds := []InputPrediction{}

	collection := c.Session.Collection(c.Name())

	err := collection.Find(as...).All(&preds)
	if err != nil {
		return nil, err
	}
	return preds, nil
}

func (m *InputPredictionCollection) Close() error {
	return nil
}
