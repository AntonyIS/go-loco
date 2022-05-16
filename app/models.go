package app

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Locomotive struct {
	LocoID      string `json:"loco_id" dynamodbav:"loco_id"`
	LocoName    string `json:"loco_name" dynamodbav:"loco_name"`
	ImageURL    string `json:"image_url" dynamodbav:"image_url" validate:"empty=false & format=url"`
	Description string `json:"description" dynamodbav:"description"`
	CreatedAT   int64  `json:"created_at" dynamodbav:"created_at"`
}

func (loco Locomotive) GetKey() map[string]types.AttributeValue {
	loco_id, err := attributevalue.Marshal(loco.LocoID)
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"loco_id": loco_id}
}
