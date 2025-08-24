package utils

import (
	"net/http"

	"golang.org/x/sync/errgroup"
)

type MyError struct {
	Msg        string
	StatusCode int
}

func (e MyError) Error() string {
	return e.Msg
}

func WaitErrGroup(g *errgroup.Group) *MyError {
	if err := g.Wait(); err != nil {
		if myErr, ok := err.(*MyError); ok {
			return myErr
		}
		return &MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}
