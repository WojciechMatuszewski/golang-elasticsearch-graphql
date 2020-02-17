package testing

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/awslabs/goformation/v4"
	cfDynamo "github.com/awslabs/goformation/v4/cloudformation/dynamodb"
)

const localDynamoAddr = "http://localhost:8000"

// GetLocalDynamo returns dynamoDB which uses local session
func GetLocalDynamo() *dynamodb.DynamoDB {
	return dynamodb.New(LocalSession(), &aws.Config{Endpoint: aws.String(localDynamoAddr)})
}

// SetupDynamoTest creates tables in DB from serverless.yml.
func SetupDynamoTest(t *testing.T, prefix string) *dynamodb.DynamoDB {
	t.Helper()

	db := GetLocalDynamo()
	tables := getDynamoTablesSchema()

	for _, table := range tables {
		resolvedTbName := table.TableName
		if strings.HasPrefix(resolvedTbName, "${") {
			parsed := strings.Split(resolvedTbName, ".")[1]
			resolvedTbName = strings.Replace(parsed, "}", "", -1)
		}

		_, err := db.CreateTable(&dynamodb.CreateTableInput{
			AttributeDefinitions: toDynamoAttributeDefs(table.AttributeDefinitions),
			KeySchema:            toDynamoKeySchema(table.KeySchema),
			BillingMode:          aws.String(table.BillingMode),
			TableName:            aws.String(prefix + resolvedTbName),
		})
		if err != nil {

			t.Fatalf(err.Error())
		}
	}

	return db
}

func getDynamoTablesSchema() map[string]*cfDynamo.Table {
	filePath, err := GetFullPath("/resources/dynamo.yml")
	if err != nil {
		panic(err.Error())
	}

	tmpl, err := goformation.Open(filePath)
	if err != nil {
		panic(err.Error())
	}

	return tmpl.GetAllDynamoDBTableResources()
}

func toDynamoAttributeDefs(defs []cfDynamo.Table_AttributeDefinition) []*dynamodb.AttributeDefinition {
	dynamoDefs := make([]*dynamodb.AttributeDefinition, len(defs))

	for i, cfDef := range defs {
		dynamoDefs[i] = &dynamodb.AttributeDefinition{
			AttributeName: aws.String(cfDef.AttributeName),
			AttributeType: aws.String(cfDef.AttributeType),
		}
	}

	return dynamoDefs
}

func toDynamoKeySchema(defs []cfDynamo.Table_KeySchema) []*dynamodb.KeySchemaElement {
	dynamoDefs := make([]*dynamodb.KeySchemaElement, len(defs))

	for i, cfDef := range defs {
		dynamoDefs[i] = &dynamodb.KeySchemaElement{
			AttributeName: aws.String(cfDef.AttributeName),
			KeyType:       aws.String(cfDef.KeyType),
		}
	}

	return dynamoDefs
}
