package db

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)


func StructToBSON(obj interface{}, logger *zap.Logger) ([]byte, error) {
	data, err := bson.Marshal(obj)
	if err != nil {
		logger.Error(fmt.Sprintf("%s : %v", "[DB] [Utils] [StructToBSON]", err))
		return nil, err
	}

	return data, nil
}
