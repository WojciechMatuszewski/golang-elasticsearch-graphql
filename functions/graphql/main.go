package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"elastic-search/pkg/env"
	"elastic-search/pkg/todo"
	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

func main() {
	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	wd, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}

	schemaF, err := os.Open(path.Join(wd, env.GRAPHQL_SCHEMA_FILE))
	if err != nil {
		panic(err.Error())
	}
	defer schemaF.Close()

	schemaB, err := ioutil.ReadAll(schemaF)
	if err != nil {
		panic(err.Error())
	}

	store := todo.NewStore(os.Getenv(env.TODO_TABLE), db)

	rootDeps := deps{store: store}

	schema := graphql.MustParseSchema(string(schemaB), &RootResolver{deps: &rootDeps})

	http.Handle("/graphql", &relay.Handler{Schema: schema})
	log.Fatal(gateway.ListenAndServe(":3000", nil))
}
