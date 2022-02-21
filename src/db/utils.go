package db

import (
	"fmt"

	"github.com/fitant/xbin-api/src/utils"
	"go.mongodb.org/mongo-driver/bson"
)


func StructToBSON(obj interface{}) ([]byte, error) {
	data, err := bson.Marshal(obj)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("%s : %v", "[DB] [Utils] [StructToBSON]", err))
		return nil, err
	}

	return data, nil
}
