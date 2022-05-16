package dynamodb

import (
	"fmt"
	"log"

	"github.com/AntonyIS/go-loco/app"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

type dynamodbRepository struct {
	Client    *dynamodb.DynamoDB
	Tablename string
}

func newDynamoDBClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create DynamoDB client
	client := dynamodb.New(sess)
	return client

}

func NewDynamoDBReposistory(tablename string) (app.LocomotiveRepository, error) {
	repo := &dynamodbRepository{
		Tablename: tablename,
	}
	repo.Client = newDynamoDBClient()
	return repo, nil
}

func (d *dynamodbRepository) CreateLoco(loco *app.Locomotive) (*app.Locomotive, error) {
	av, err := dynamodbattribute.MarshalMap(*loco)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Locomotive.CreateLoco")
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(d.Tablename),
	}
	_, err = d.Client.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
	return loco, nil
}

func (d *dynamodbRepository) GetLoco(loco_id string) (*app.Locomotive, error) {
	result, err := d.Client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(d.Tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"loco_id": {
				S: aws.String(loco_id),
			},
		},
	})

	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}
	if result.Item == nil {
		return nil, errors.Wrap(app.ErrorLocomotiveNotFound, "repository.Locomotive.GetLoco")
	}
	loco := app.Locomotive{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &loco)
	if err != nil {
		return nil, errors.Wrap(app.ErrorInternalServer, "repository.Locomotive.GetLoco")
	}
	return &loco, err
}

func (d *dynamodbRepository) UpdateLoco(loco *app.Locomotive) (*app.Locomotive, error) {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":n": {
				S: aws.String(loco.LocoName),
			},
			":i": {
				S: aws.String(loco.ImageURL),
			},
			":d": {
				S: aws.String(loco.Description),
			},
		},
		TableName: aws.String(d.Tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"loco_id": {
				S: aws.String(loco.LocoID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set loco_name = :n,image_url = :i,description = :d"),
	}
	te, err := d.Client.UpdateItem(input)
	fmt.Println(te)
	if err != nil {
		log.Fatalf("Got error calling UpdateItem: %s", err)
	}
	return loco, nil
}

func (d *dynamodbRepository) DeleteLoco(loco_id string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"loco_id": {
				S: aws.String(loco_id),
			},
		},
		TableName: aws.String(d.Tablename),
	}
	_, err := d.Client.DeleteItem(input)
	if err != nil {
		log.Fatalf("Got error calling DeleteItem: %s", err)
	}

	return nil

}
func (d *dynamodbRepository) GetAllLoco() (*[]app.Locomotive, error) {
	locos := &[]app.Locomotive{}
	params := &dynamodb.ScanInput{
		TableName: aws.String(d.Tablename),
	}
	// Make the DynamoDB Query API call
	result, err := d.Client.Scan(params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}
	dynamodbattribute.UnmarshalListOfMaps(result.Items, locos)

	return locos, nil
}
