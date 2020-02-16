package main

import "elastic-search/pkg/todo"

type RootResolver struct {
	deps *deps
}

// StoreIface is a interface which represents the store
type StoreIface interface {
	Save(todo todo.Todo) error
	GetByID(ID string) (todo.Todo, error)
}


type deps struct {
	store StoreIface
}





