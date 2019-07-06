package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hikaru7719/lambda-bookshelf/finder"
	"github.com/hikaru7719/lambda-bookshelf/register"
)

func BookRegister() {
	book, err := finder.Find("9784797394481")
	if err == nil {
		register.Upload(book)
	}
}

func main() {
	lambda.Start(BookRegister)
}
