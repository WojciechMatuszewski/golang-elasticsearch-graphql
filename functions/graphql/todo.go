package main

import (
	"context"

	"elastic-search/pkg/todo"
	"github.com/graph-gophers/graphql-go"
	"github.com/rs/xid"
)

// CreateTodoInput represents a structure of argument passed to CreateTodo mutation.
type CreateTodoInput struct {
	Content string `json:"content"`
}

// CreateTodo creates Todo using todo Store.
func (r *RootResolver) CreateTodo(args struct{ Input CreateTodoInput }) (*TodoResolver, error) {
	td := todo.Todo{
		ID:      xid.New().String(),
		Content: args.Input.Content,
	}

	err := r.deps.store.Save(td)
	if err != nil {
		return nil, err
	}

	return &TodoResolver{todo: td}, nil
}

// GetTodo gets the Todo by ID using todo Store.
func (r *RootResolver) GetTodo(args struct{ ID graphql.ID }) (*TodoResolver, error) {
	td, err := r.deps.store.GetByID(string(args.ID))
	if err != nil {
		return nil, err
	}

	return &TodoResolver{todo: td}, nil
}

type SearchTodoInput struct {
	Query string
}

func (r *RootResolver) SearchTodo(ctx context.Context, args SearchTodoInput) ([]*TodoResolver, error) {
	tds, err := r.deps.esService.Search(ctx, args.Query)
	if err != nil {
		return nil, err
	}

	resolvers := make([]*TodoResolver, len(tds))
	for i, td := range tds {
		resolvers[i] = &TodoResolver{todo: td}
	}

	return resolvers, nil
}

// TodoResolver resolves Todo-related graphql fields
type TodoResolver struct {
	todo todo.Todo
}

func (tdR *TodoResolver) ID() graphql.ID {
	return graphql.ID(tdR.todo.ID)
}

func (tdR *TodoResolver) Content() string {
	return tdR.todo.Content
}
