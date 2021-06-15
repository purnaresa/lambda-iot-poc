package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Temperature struct {
	NodeID    string `json:"node_id"`
	Timestamp string `json:"timestamp"`
	Value     string `json:"value"`
}

type CommandInput struct {
	NodeID      string `json:"node_id"`
	Temperature string `json:"temperature"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	svc := dynamodb.New(session.New())
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String("101"),
			},
		},
		KeyConditionExpression: aws.String("NodeId = :v1"),
		ProjectionExpression:   aws.String("RecordTime"),
		TableName:              aws.String("Temp"),
	}

	result, err := svc.Query(input)
	if err != nil {
		log.Println(err)
	}
	output := []Temperature{}
	for _, v := range result.Items {
		o := Temperature{}
		log.Println(o)
		err = dynamodbattribute.UnmarshalMap(v, &o)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal data, %v", err))
		}
		output = append(output, o)
	}

	outputByte, _ := json.Marshal(&output)
	outputJson := string(outputByte)

	return events.APIGatewayProxyResponse{
		Body:       outputJson,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
