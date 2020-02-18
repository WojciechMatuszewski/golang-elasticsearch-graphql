package main

import (
	"context"

	"elastic-search/pkg/todo"
)

type RootResolver struct {
	deps *deps
}

// StoreIface is an interface which represents the store
type StoreIface interface {
	Save(todo todo.Todo) error
	GetByID(ID string) (todo.Todo, error)
}

// ElasticSearchServiceIface is an interface which represents the Elastic Search service
type ElasticSearchServiceIface interface {
	Index(ctx context.Context, td todo.Todo) error
	Search(ctx context.Context, query string) ([]todo.Todo, error)
}

type deps struct {
	store     StoreIface
	esService ElasticSearchServiceIface
}
