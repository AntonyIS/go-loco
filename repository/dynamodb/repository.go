package dynamodb

import (
	"context"

	"github.com/AntonyIS/go-loco/app"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/pkg/errors"
)

type dynamodbRepository struct {
	Client    *dynamodb.Client
	Tablename string
}

func NewDynamoDBReposistory(tablename string) (app.LocomotiveRepository, error) {
	repo := &dynamodbRepository{
		Tablename: tablename,
	}
	return repo, nil
}

func (d *dynamodbRepository) CreateLoco(loco *app.Locomotive) (*app.Locomotive, error) {
	item, err := attributevalue.MarshalMap(*loco)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.CreateLoco")
	}
	_, err = d.Client.PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			TableName: aws.String(d.Tablename),
			Item:      item,
		})
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.CreateLoco")

	}

	return loco, nil
}

func (d *dynamodbRepository) GetLoco(loco_id string) (*app.Locomotive, error) {
	loco := &app.Locomotive{LocoID: loco_id}
	response, err := d.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       loco.GetKey(),
		TableName: aws.String(d.Tablename),
	})
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.GetLoco")
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &response)
		if err != nil {
			return nil, errors.Wrap(err, "repository.Redirect.GetLoco")
		}
	}
	return loco, err
}

func (d *dynamodbRepository) UpdateLoco(loco *app.Locomotive) (*app.Locomotive, error) {
	var response *dynamodb.UpdateItemOutput
	var attributeMap map[string]map[string]interface{}
	update := expression.Set(expression.Name("name"), expression.Value(loco.Name))
	update.Set(expression.Name("image_url"), expression.Value(loco.ImageURL))
	update.Set(expression.Name("description"), expression.Value(loco.Description))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.GetLoco")
	} else {
		response, err = d.Client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
			TableName: aws.String(d.Tablename),
			Key:       loco.GetKey(),
			// ExpressionAttributeNames:  expr.Names(),
			// ExpressionAttributeValues: expr.Values(),
			UpdateExpression: expr.Update(),
			ReturnValues:     types.ReturnValueUpdatedNew,
		})
		if err != nil {
			return nil, errors.Wrap(err, "repository.Redirect.GetLoco")
		} else {
			err = attributevalue.UnmarshalMap(response.Attributes, &attributeMap)
			if err != nil {
				return nil, errors.Wrap(err, "repository.Redirect.GetLoco")
			}
		}
	}
	return loco, err
}

func (d *dynamodbRepository) DeleteLoco(loco *app.Locomotive) error {
	_, err := d.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(d.Tablename), Key: loco.GetKey(),
	})
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}
