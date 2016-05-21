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

const table = "counters"

type message struct {
	SerialNumber   string `json:"serialNumber"`
	BatteryVoltage string `json:"batteryVoltage"`
	ClickType      string `json:"clickType"`
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		var m message

		if err := json.Unmarshal(event, &m); err != nil {
			return nil, err
		}

		id := os.Getenv("counterId")
		region := os.Getenv("region")
		client := createDynamoClient(region)

		return incrementCounter(client, id)
	})
}

func createDynamoClient(region string) *dynamodb.DynamoDB {
	session := session.New()
	config := aws.NewConfig().WithRegion(region)
	dynamo := dynamodb.New(session, config)
	return dynamo
}

func incrementCounter(dynamo *dynamodb.DynamoDB, id string) (int64, error) {
	params := &dynamodb.UpdateItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"counter_id": {
				S: aws.String(id),
			},
		},
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{
			"count": {
				Action: aws.String("ADD"),
				Value: &dynamodb.AttributeValue{
					N: aws.String("1"),
				},
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	resp, err := dynamo.UpdateItem(params)
	if err != nil {
		return 0, err
	}

	countAttr := resp.Attributes["count"]
	countValue := aws.StringValue(countAttr.N)
	countInt, err := strconv.ParseInt(countValue, 10, 64)

	return countInt, err
}
