package base

type Repository interface {
	GetConnection() (T any)
	GetModelName() string
}
