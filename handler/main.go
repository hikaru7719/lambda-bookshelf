package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	l "github.com/aws/aws-sdk-go/service/lambda"
	"net/http"
	"os"
	"strings"
)

type Find struct {
	Word        string `json:"word"`
	ChannelName string `json:"channel_name"`
}

func SlackHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(request.Body), &jsonMap)

	eventType := jsonMap["type"].(string)

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
	}

	switch eventType {
	case "url_verification":
		challenge := jsonMap["challenge"].(string)
		response.Body = challenge
		return response, nil
	case "event_callback":
		event := jsonMap["event"].(map[string]interface{})
		eventTypeString := event["type"].(string)
		if eventTypeString == "app_mention" {
			eventText := event["text"].(string)
			stringReader := strings.NewReader(eventText)
			scanner := bufio.NewScanner(stringReader)
			scanner.Scan()
			scanner.Scan()
			text := scanner.Text()
			splitSlice := strings.Split(text, ":")
			if len(splitSlice) < 1 {
				response.StatusCode = http.StatusBadRequest
				return response, errors.New("search word should include `:`")
			}
			switch splitSlice[0] {
			case "search":
				channelName := event["channel"].(string)
				find := Find{Word: splitSlice[1], ChannelName: channelName}
				jsonBytes, _ := json.Marshal(find)
				newSession, _ := session.NewSession(&aws.Config{
					Region: aws.String("ap-northeast-1")},
				)
				lambdaPublisher := l.New(newSession)
				input := &l.InvokeInput{
					FunctionName:   aws.String(os.Getenv("LAMBDA_FUNC_NAME")),
					Payload:        jsonBytes,
					InvocationType: aws.String("Event"),
				}

				lambdaPublisher.Invoke(input)
				return response, nil
			}
		}
	}

	return response, nil
}

func main() {
	lambda.Start(SlackHandler)
}
