package main

import (
	testing2 "elastic-search/pkg/testing"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func deleteTables(db *dynamodb.DynamoDB) error {
	tables, err := db.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return err
	}

	for _, tn := range tables.TableNames {
		_, err = db.DeleteTable(&dynamodb.DeleteTableInput{TableName: tn})
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	db := testing2.GetLocalDynamo()
	err := deleteTables(db)
	if err != nil {
		panic(err.Error())
	}
}
