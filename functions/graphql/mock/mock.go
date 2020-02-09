//go:generate mockgen -destination=graphql.go  -package=mock elastic-search/pkg/todo StoreIface
package mock
