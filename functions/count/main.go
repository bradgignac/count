package main

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/apex/go-apex"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
)

type message struct {
	ClickType string `json:"clickType"`
}

var region = os.Getenv("region")
var endpoint = os.Getenv("endpoint")
var iotCtrl *iot.IoT
var iotData *iotdataplane.IoTDataPlane
var dynamo *dynamodb.DynamoDB

func init() {
	iotCtrl = createIotClient()
	iotData = createIotDataplaneClient()
	dynamo = createDynamoClient()
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		var m message

		if err := json.Unmarshal(event, &m); err != nil {
			return nil, err
		}

		id := os.Getenv("counterId")

		switch m.ClickType {
		case "LONG":
			return createCounter(id)
		case "SINGLE":
			return incrementCounter(id)
		case "DOUBLE":
			return resetCounter(id)
		}

		return nil, nil
	})
}

func createIotClient() *iot.IoT {
	session := session.New()
	config := aws.NewConfig().WithRegion(region)
	return iot.New(session, config)
}

func createIotDataplaneClient() *iotdataplane.IoTDataPlane {
	session := session.New()
	config := aws.NewConfig().
		WithRegion(region).
		WithEndpoint(endpoint)
	return iotdataplane.New(session, config)
}

func createDynamoClient() *dynamodb.DynamoDB {
	session := session.New()
	config := aws.NewConfig().WithRegion(region)
	return dynamodb.New(session, config)
}

func createCounter(id string) (string, error) {
	params := &iot.CreateThingInput{ThingName: aws.String(id)}
	resp, err := iotCtrl.CreateThing(params)
	if err != nil {
		return "", err
	}

	return aws.StringValue(resp.ThingArn), nil
}

func incrementCounter(id string) (int64, error) {
	count, err := updateCounter(id, map[string]*dynamodb.AttributeValueUpdate{
		"count": {
			Action: aws.String("ADD"),
			Value: &dynamodb.AttributeValue{
				N: aws.String("1"),
			},
		},
	})
	if err != nil {
		return 0, err
	}

	return updateShadow(id, count)
}

func resetCounter(id string) (int64, error) {
	count, err := updateCounter(id, map[string]*dynamodb.AttributeValueUpdate{
		"count": {
			Action: aws.String("PUT"),
			Value: &dynamodb.AttributeValue{
				N: aws.String("0"),
			},
		},
	})
	if err != nil {
		return 0, err
	}

	return updateShadow(id, count)
}

func updateCounter(id string, update map[string]*dynamodb.AttributeValueUpdate) (int64, error) {
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

func updateShadow(id string, count int64) (int64, error) {
	shadow := &map[string]interface{}{
		"state": &map[string]interface{}{
			"desired": &map[string]interface{}{"count": count},
		},
	}
	payload, err := json.Marshal(shadow)
	if err != nil {
		return 0, err
	}

	params := &iotdataplane.UpdateThingShadowInput{
		ThingName: aws.String(id),
		Payload:   payload,
	}
	if _, err = iotData.UpdateThingShadow(params); err != nil {
		return 0, err
	}

	return count, nil
}
