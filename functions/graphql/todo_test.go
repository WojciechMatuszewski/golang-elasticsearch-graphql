package main

import (
	"io/ioutil"
	"os"
	"testing"

	"elastic-search/functions/graphql/mock"
	"elastic-search/pkg/env"
	testing2 "elastic-search/pkg/testing"
	"elastic-search/pkg/todo"
	"github.com/golang/mock/gomock"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/gqltesting"
	"github.com/rs/xid"
)

var schemaB []byte

func TestMain(m *testing.M) {
	schemaP, err := testing2.GetFullPath("/" + env.GRAPHQL_SCHEMA_FILE)
	if err != nil {
		panic(err.Error())
	}

	schemaF, err := os.Open(schemaP)
	if err != nil {
		panic(err.Error())
	}
	defer schemaF.Close()

	b, err := ioutil.ReadAll(schemaF)
	if err != nil {
		panic(err.Error())
	}

	schemaB = b

	m.Run()
}

func TestRootResolver_GetTodo(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		store := mock.NewMockStoreIface(ctrl)
		resolverDeps := deps{store: store}
		schema := graphql.MustParseSchema(string(schemaB), &RootResolver{deps: &resolverDeps})
		todoID := xid.New().String()
		outTodo := todo.Todo{
			ID:      todoID,
			Content: "content",
		}

		store.EXPECT().GetByID(todoID).Return(outTodo, nil)
		gqltesting.RunTest(t, &gqltesting.Test{
			Context: nil,
			Schema:  schema,
			Query: `
				query q($ID: ID!){
					getTodo(ID: $ID){
						content
						ID
					}
				}
			`,
			Variables: map[string]interface{}{
				"ID": todoID,
			},
			ExpectedResult: `{
				"getTodo":
					{
						"content":"content",
						"ID":"` + outTodo.ID + `"
					}
				}`,
		})
	})

}
