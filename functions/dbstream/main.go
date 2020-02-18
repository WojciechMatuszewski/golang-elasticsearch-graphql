package main

import (
	"os"

	"elastic-search/pkg/env"
	"elastic-search/platform/elasticsearch"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {
	sess := session.Must(session.NewSession())

	indexer, err := elasticsearch.NewService(sess, os.Getenv(env.ELASTIC_SEARCH_ENDPOINT))
	if err != nil {
		panic(err.Error())
	}

	handler := NewHandler(indexer)

	lambda.Start(handler)
}
