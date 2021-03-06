package todo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// Store represents TodoStore
type Store struct {
	tableName string
	db        *dynamodb.DynamoDB
}

// NewStore creates TodoStore
func NewStore(tableName string, db *dynamodb.DynamoDB) *Store {
	return &Store{
		tableName: tableName,
		db:        db,
	}
}

// TODO expression is not viable when doing `getItems` simple operations

// GetByID gets given Todo by ID
func (s *Store) GetByID(ID string) (Todo, error) {
	out, err := s.db.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(ID),
			},
		},
		TableName: aws.String(s.tableName),
	})

	if err != nil {
		return Todo{}, err
	}

	var outTodo Todo
	err = dynamodbattribute.UnmarshalMap(out.Item, &outTodo)
	if err != nil {
		return Todo{}, err
	}

	return outTodo, nil
}

// Save saves given todo to database
func (s *Store) Save(todo Todo) error {
	in, err := dynamodbattribute.MarshalMap(todo)
	if err != nil {
		return err
	}

	_, err = s.db.PutItem(&dynamodb.PutItemInput{
		Item:      in,
		TableName: aws.String(s.tableName),
	})

	if err != nil {
		return err
	}

	return nil
}

// Remove removes given todo from the database
func (s *Store) Remove(ID string) error {
	cond := expression.AttributeExists(expression.Name("ID"))
	expr, err := expression.NewBuilder().WithCondition(cond).Build()
	if err != nil {
		return err
	}

	_, err = s.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(ID),
			},
		},
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return err
	}
	return nil
}
