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

ERROR 1062 (23000): Duplicate entry '1-2' for key 'followers.unique_follow'


 */