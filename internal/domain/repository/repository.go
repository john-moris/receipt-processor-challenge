package repository

import "errors"

var (
	ErrItemNotFound        = errors.New("item with the given id does not exist")
	ErrItemStillInProgress = errors.New("item with the given id is still processing")
)

type Repository interface {
	Get(id string) (int, error)
	Start() string
	Finish(id string, score int)
}
