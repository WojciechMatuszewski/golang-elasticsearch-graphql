package main

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"elastic-search/functions/graphql/mock"
	"elastic-search/pkg/env"
	testing2 "elastic-search/pkg/testing"
	"elastic-search/pkg/todo"
	"github.com/golang/mock/gomock"
	"github.com/graph-gophers/graphql-go"
	graphqlerrors "github.com/graph-gophers/graphql-go/errors"
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

func TestRootResolver_CreateTodo(t *testing.T) {
	t.Parallel()

	t.Run("failure on saving the todo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		store := mock.NewMockStoreIface(ctrl)
		resolverDeps := deps{store: store}
		schema := graphql.MustParseSchema(string(schemaB), &RootResolver{deps: &resolverDeps})

		in := CreateTodoInput{Content: "content"}
		store.EXPECT().Save(gomock.Any()).Return(errors.New("boom"))
		gqltesting.RunTest(t, &gqltesting.Test{
			Schema: schema,
			Query: `
	mutation createTodoMutation($input: CreateTodoInput!){
		createTodo(input: $input){
			content
		}
	}
`,
			Variables: map[string]interface{}{
				"input": structToMap(&in),
			},
			ExpectedResult: `null`,
			ExpectedErrors: []*graphqlerrors.QueryError{
				{
					Message:       "boom",
					Path:          []interface{}{"createTodo"},
					ResolverError: errors.New("boom"),
				},
			},
		})
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		store := mock.NewMockStoreIface(ctrl)
		resolverDeps := deps{store: store}
		schema := graphql.MustParseSchema(string(schemaB), &RootResolver{deps: &resolverDeps})

		// what todo??
		in := CreateTodoInput{Content: "content"}

		store.EXPECT().Save(gomock.Any()).Return(nil)
		gqltesting.RunTest(t, &gqltesting.Test{
			Schema: schema,
			Query: `
	mutation createTodoMutation($input: CreateTodoInput!){
		createTodo(input: $input){
			content
		}
	}
`,
			Variables: map[string]interface{}{
				"input": structToMap(&in),
			},
			ExpectedResult: `{"createTodo":{"content":"content"}}`,
			ExpectedErrors: nil,
		})
	})
}

func TestRootResolver_GetTodo(t *testing.T) {
	t.Parallel()

	t.Run("failure on getByID error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		store := mock.NewMockStoreIface(ctrl)
		resolverDeps := deps{store: store}
		schema := graphql.MustParseSchema(string(schemaB), &RootResolver{deps: &resolverDeps})
		todoID := xid.New().String()

		store.EXPECT().GetByID(todoID).Return(todo.Todo{}, errors.New("boom"))
		gqltesting.RunTest(t, &gqltesting.Test{
			Schema: schema,
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
			ExpectedResult: `{"getTodo":null}`,
			ExpectedErrors: []*graphqlerrors.QueryError{
				{
					Message:       "boom",
					Path:          []interface{}{"getTodo"},
					ResolverError: errors.New("boom"),
				},
			},
		})
	})

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

func TestRootResolver_Search(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("success multiple found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		esService := mock.NewMockElasticSearchServiceIface(ctrl)
		resolverDeps := deps{esService: esService}
		schema := graphql.MustParseSchema(string(schemaB), &RootResolver{deps: &resolverDeps})

		outTds := []todo.Todo{
			{
				ID:      xid.New().String(),
				Content: "so",
			},
			{
				ID:      xid.New().String(),
				Content: "con",
			},
		}
		in := SearchTodoInput{Query: "some content"}

		esService.EXPECT().Search(gomock.Any(), in.Query).Return(outTds, nil)
		gqltesting.RunTest(t, &gqltesting.Test{
			Context: ctx,
			Schema:  schema,
			Query: `
query searchForTodos($query: String!){
	searchTodo(query: $query) {
		content
}
}
`,
			Variables: map[string]interface{}{
				"query": in.Query,
			},
			ExpectedResult: `{"searchTodo":[{"content":"so"},{"content":"con"}]}`,
			ExpectedErrors: nil,
		})
	})

	t.Run("success 0 found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		esService := mock.NewMockElasticSearchServiceIface(ctrl)
		rootDeps := deps{esService: esService}
		schema := graphql.MustParseSchema(string(schemaB), &RootResolver{deps: &rootDeps})

		in := SearchTodoInput{Query: "some content"}

		esService.EXPECT().Search(gomock.Any(), in.Query).Return([]todo.Todo{}, nil)
		gqltesting.RunTest(t, &gqltesting.Test{
			Context: ctx,
			Schema:  schema,
			Query: `
query searchForTodos($query: String!){
	searchTodo(query: $query) {
		content
}
}
`,
			Variables: map[string]interface{}{
				"query": in.Query,
			},
			ExpectedResult: `{"searchTodo": []}`,
			ExpectedErrors: nil,
		})
	})

	t.Run("failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		esService := mock.NewMockElasticSearchServiceIface(ctrl)
		rootDeps := deps{esService: esService}
		schema := graphql.MustParseSchema(string(schemaB), &RootResolver{deps: &rootDeps})

		in := SearchTodoInput{Query: "some content"}

		esService.EXPECT().Search(gomock.Any(), in.Query).Return([]todo.Todo{}, errors.New("boom"))
		gqltesting.RunTest(t, &gqltesting.Test{
			Context: ctx,
			Schema:  schema,
			Query: `
query searchForTodos($query: String!){
	searchTodo(query: $query) {
		content
}
}
`,
			Variables: map[string]interface{}{
				"query": in.Query,
			},
			ExpectedResult: `null`,
			ExpectedErrors: []*graphqlerrors.QueryError{
				{
					Message:       "boom",
					Path:          []interface{}{"searchTodo"},
					ResolverError: errors.New("boom"),
				},
			},
		})
	})
}

func structToMap(i interface{}) map[string]interface{} {
	out := map[string]interface{}{}
	iVal := reflect.ValueOf(i).Elem()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		var v string
		switch f.Interface().(type) {
		case string:
			v = f.String()
			out[v] = f.Interface()
			break
		default:
			return out
		}

	}
	return out
}
