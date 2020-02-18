package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"elastic-search/functions/dbstream/mock"
	testing2 "elastic-search/pkg/testing"
	"elastic-search/pkg/todo"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	"github.com/tj/assert"
)

func TestNewHandler(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("success single", func(t *testing.T) {
		evt := getLambdaPayload(t, "single.json")
		ctrl := gomock.NewController(t)
		indexer := mock.NewMockIndexer(ctrl)
		handler := NewHandler(indexer)

		var inTd todo.Todo
		err := unmarshalStreamImage(evt.Records[0].Change.NewImage, &inTd)
		if err != nil {
			t.Fatalf(err.Error())
		}

		indexer.EXPECT().Index(ctx, inTd).Return(nil)

		err = handler(ctx, evt)
		assert.NoError(t, err)
	})

	t.Run("success multiple", func(t *testing.T) {
		evt := getLambdaPayload(t, "multiple.json")
		ctrl := gomock.NewController(t)
		indexer := mock.NewMockIndexer(ctrl)
		handler := NewHandler(indexer)

		for _, dbEvt := range evt.Records {
			var inTd todo.Todo
			err := unmarshalStreamImage(dbEvt.Change.NewImage, &inTd)
			if err != nil {
				t.Fatalf(err.Error())
			}

			indexer.EXPECT().Index(ctx, inTd).Return(nil)
		}

		err := handler(ctx, evt)
		assert.NoError(t, err)
	})

	t.Run("indexer failure", func(t *testing.T) {
		evt := getLambdaPayload(t, "multiple.json")
		ctrl := gomock.NewController(t)
		indexer := mock.NewMockIndexer(ctrl)
		handler := NewHandler(indexer)

		var inTd todo.Todo
		err := unmarshalStreamImage(evt.Records[0].Change.NewImage, &inTd)
		if err != nil {
			t.Fatalf(err.Error())
		}

		indexer.EXPECT().Index(ctx, inTd).Return(errors.New("boom"))

		err = handler(ctx, evt)
		assert.NoError(t, err)
	})
}

func getLambdaPayload(t *testing.T, fName string) events.DynamoDBEvent {
	fPath, err := testing2.GetFullPath("/functions/dbstream/testdata/" + fName)
	if err != nil {
		t.Fatalf(err.Error())
	}

	f, err := os.Open(fPath)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fB, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf(err.Error())
	}

	var e events.DynamoDBEvent
	err = json.Unmarshal(fB, &e)
	if err != nil {
		t.Fatalf(err.Error())
	}

	return e
}
