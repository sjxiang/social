package store

import "errors"

var (
	ErrNotFound          = errors.New("resource not found")
	ErrAlreadyExists     = errors.New("resource already exists")
	ErrConflict          = errors.New("resource already exists")
	
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)


/*

增, already exists
	
删改, 0 rows affected

查, not found

 */