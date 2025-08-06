package domain

type ModelProvider[T any] interface {
	ListModel(subType string) ([]T, error)
}