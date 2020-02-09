package main

import "elastic-search/pkg/todo"

type RootResolver struct {
	deps *deps
}

type deps struct {
	store todo.StoreIface
}
