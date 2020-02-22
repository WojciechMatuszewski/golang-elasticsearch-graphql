package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"reflect"

	"elastic-search/pkg/todo"

	"github.com/apex/log"
	apexJSON "github.com/apex/log/handlers/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodbstreams"
)

// Handler represents a lambda handler function
type Handler func(ctx context.Context, e events.DynamoDBEvent) error

// Indexer indexs a todo to elastic search service
type IndexerRemover interface {
	Index(ctx context.Context, td todo.Todo) error
	Remove(ctx context.Context, ID string) error
}

// NewHandler is the function responsible for creating dbstream lambda handler
func NewHandler(indexerRemover IndexerRemover) Handler {
	log.SetHandler(apexJSON.New(os.Stdout))
	return func(ctx context.Context, e events.DynamoDBEvent) error {
		for _, evt := range e.Records {
			log.WithField("event", evt).Info("incoming event")
			switch evt.EventName {
			case dynamodbstreams.OperationTypeRemove:
				log.Info(dynamodbstreams.OperationTypeRemove)

				var td todo.Todo
				err := unmarshalStreamImage(evt.Change.OldImage, &td)
				if err != nil {
					log.WithError(err).Error("error while unmarshaling stream image")
					return nil
				}

				err = indexerRemover.Remove(ctx, td.ID)
				if err != nil {
					log.WithError(err).Error("error while removing from es")
					return nil
				}

			case dynamodbstreams.OperationTypeInsert:
				log.Info(dynamodbstreams.OperationTypeInsert)

				var td todo.Todo
				err := unmarshalStreamImage(evt.Change.NewImage, &td)

				if err != nil {
					log.WithError(err).Error("error while unmarshalling stream image")
					return nil
				}

				err = indexerRemover.Index(ctx, td)
				if err != nil {
					log.WithError(err).Error("error while indexing to elastic search")
					return nil
				}
			}
		}
		return nil
	}
}

func unmarshalStreamImage(image map[string]events.DynamoDBAttributeValue, out interface{}) error {
	outK := reflect.ValueOf(out).Kind()
	if outK != reflect.Ptr {
		return errors.New("not pointer")
	}

	attrsMap := make(map[string]*dynamodb.AttributeValue, len(image))

	for k, v := range image {
		vB, err := v.MarshalJSON()
		if err != nil {
			return err
		}

		var attr dynamodb.AttributeValue
		err = json.Unmarshal(vB, &attr)
		if err != nil {
			return err
		}

		attrsMap[k] = &attr
	}

	return dynamodbattribute.UnmarshalMap(attrsMap, out)
}
