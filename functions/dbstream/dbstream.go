package main

import (
	"encoding/json"
	"errors"
	"reflect"

	"elastic-search/pkg/todo"
	"github.com/apex/log"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodbstreams"
)

type Handler func(e events.DynamoDBEvent) error

func NewHandler() Handler {
	return func(e events.DynamoDBEvent) error {
		for _, evt := range e.Records {
			switch evt.EventName {
			case dynamodbstreams.OperationTypeInsert:

				var td todo.Todo
				err := unmarshalStreamImage(evt.Change.NewImage, &td)

				if err != nil {
					log.WithError(err).Error("error while unmarshalling stream image")
					return nil
				}

				log.WithField("todo", td).Info("unmarshalled")
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
