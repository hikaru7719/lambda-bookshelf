package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hikaru7719/lambda-bookshelf/finder"
	"github.com/hikaru7719/lambda-bookshelf/sender"
	"os"
	"strconv"
	"strings"
)

type Find struct {
	Word        string `json:"word"`
	ChannelName string `json:"channel_name"`
}

func SendMessage(find Find) error {
	newSession, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")},
	)
	s3Finder := s3.New(newSession)
	bucket := os.Getenv("BUCKET")
	fileName := os.Getenv("FILE_NAME")
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}

	result, err := s3Finder.GetObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return err
	}

	fmt.Println(result)
	csvFinder := finder.CSV{Reader: result.Body}
	bookSlice, err := csvFinder.Find(find.Word)
	if err != nil {
		fmt.Print(err)
		return err
	}

	var sendMessage strings.Builder
	length := strconv.Itoa(len(bookSlice))
	fmt.Println(length)
	if len(bookSlice) > 0 {
		sendMessage.WriteString("本あったよ:sunglasses:\n")
		sendMessage.WriteString("検索結果/" + length + "件\n")
	} else {
		sendMessage.WriteString("残念だ...\n")
	}

	for _, book := range bookSlice {
		sendMessage.WriteString("```")
		sendMessage.WriteString(book.ToString())
		sendMessage.WriteString("```")
		sendMessage.WriteString("\n")
	}
	sender.SendMessage(find.ChannelName, sendMessage.String())
	return nil
}

func main() {
	lambda.Start(SendMessage)
}
