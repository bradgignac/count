package main

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/apex/go-apex"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type message struct {
	ClickType string `json:"clickType"`
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		var m message

		if err := json.Unmarshal(event, &m); err != nil {
			return nil, err
		}

		id := os.Getenv("counterId")
		client := createDynamoClient()

		switch m.ClickType {
		case "SINGLE":
			return incrementCounter(client, id)
		case "LONG":
			return resetCounter(client, id)
		}

		return nil, nil
	})
}

func createDynamoClient() *dynamodb.DynamoDB {
	region := os.Getenv("region")
	session := session.New()
	config := aws.NewConfig().WithRegion(region)
	dynamo := dynamodb.New(session, config)
	return dynamo
}

func incrementCounter(dynamo *dynamodb.DynamoDB, id string) (int64, error) {
	return updateCounter(dynamo, id, map[string]*dynamodb.AttributeValueUpdate{
		"count": {
			Action: aws.String("ADD"),
			Value: &dynamodb.AttributeValue{
				N: aws.String("1"),
			},
		},
	})
}

func resetCounter(dynamo *dynamodb.DynamoDB, id string) (int64, error) {
	return updateCounter(dynamo, id, map[string]*dynamodb.AttributeValueUpdate{
		"count": {
			Action: aws.String("PUT"),
			Value: &dynamodb.AttributeValue{
				N: aws.String("0"),
			},
		},
	})
}

func updateCounter(dynamo *dynamodb.DynamoDB, id string, update map[string]*dynamodb.AttributeValueUpdate) (int64, error) {
	params := &dynamodb.UpdateItemInput{
		TableName: aws.String("counters"),
		Key: map[string]*dynamodb.AttributeValue{
			"counter_id": {
				S: aws.String(id),
			},
		},
		AttributeUpdates: update,
		ReturnValues:     aws.String("UPDATED_NEW"),
	}

	resp, err := dynamo.UpdateItem(params)
	if err != nil {
		return 0, err
	}

	countAttr := resp.Attributes["count"]
	countValue := aws.StringValue(countAttr.N)
	return strconv.ParseInt(countValue, 10, 64)
}
