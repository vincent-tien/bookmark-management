package errors

import "errors"

var ErrKeyAlreadyExists = errors.New("key already exists")
var ErrUrlNotFound = errors.New("url not found")
var ErrInvalidAuth = errors.New("invalid username or password")
