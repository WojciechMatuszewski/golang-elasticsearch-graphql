package main

import (
	"testing"

	testing2 "elastic-search/pkg/testing"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rs/xid"
	"github.com/tj/assert"
)

func TestDeleteTables(t *testing.T) {
	t.Run("no tables", func(t *testing.T) {
		db := testing2.GetLocalDynamo()

		err := deleteTables(db)
		assert.NoError(t, err)
	})

	t.Run("success", func(t *testing.T) {
		prefix := xid.New().String()
		db := testing2.SetupDynamoTest(t, prefix)

		err := deleteTables(db)
		assert.NoError(t, err)

		tables, err := db.ListTables(&dynamodb.ListTablesInput{})
		if err != nil {
			t.Fatalf(err.Error())
		}

		assert.Len(t, tables.TableNames, 0)
	})
}
