package register

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hikaru7719/lambda-bookshelf/finder"
	"io/ioutil"
	"os"
)

func Upload(book *finder.Book) error {
	newSession, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")},
	)
	s3Management := s3.New(newSession)
	bucket := os.Getenv("BUCKET")
	fileName := os.Getenv("FILE_NAME")
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}

	result, err := s3Management.GetObject(input)
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
	file, _ := os.Create("/tmp/data.csv")
	buf, _ := ioutil.ReadAll(result.Body)
	file.Write(buf)
	file.WriteString(book.ConvertCSV())
	file.Close()
	readFile, _ := os.Open("/tmp/data.csv")
	output := &s3.PutObjectInput{
		Body:   readFile,
		Bucket: aws.String(bucket),
		Key:    aws.String("book.csv"),
	}
	result2, err2 := s3Management.PutObject(output)
	if err2 != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return err2
	}
	fmt.Println(result2)
	return nil
}
