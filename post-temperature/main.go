package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type CommandOutput struct {
	Message string `json:"message"`
}

type CommandInput struct {
	NodeID string `json:"node_id"`
	Data   string `json:"data"`
}

func generateOutput(message string) (output string) {
	commandOutput := CommandOutput{
		Message: "command success",
	}

	outputByte, _ := json.Marshal(&commandOutput)
	output = string(outputByte)
	return
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var command CommandInput
	jsonErr := json.Unmarshal([]byte(request.Body), &command)
	if jsonErr != nil {
		return events.APIGatewayProxyResponse{}, errors.New(jsonErr.Error())
	}
	ts := fmt.Sprintf("%d", time.Now().UTC().Unix())
	svc := dynamodb.New(session.New())
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"NodeId": {
				S: aws.String(command.NodeID),
			},
			"Data": {
				S: aws.String(command.Data),
			},
			"RecordTime": {
				S: aws.String(ts),
			},
		},
		TableName: aws.String("Temp"),
	}

	_, err := svc.PutItem(input)
	if err != nil {
		log.Println(err)
	}

	return events.APIGatewayProxyResponse{
		Body:       generateOutput("Command successfull"),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
