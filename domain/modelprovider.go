package domain

type IModelProvider[T any] interface {
	ListModel(subType string, provider string) ([]T, error)
}