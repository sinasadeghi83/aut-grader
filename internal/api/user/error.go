package user

import "errors"

var ErrInvalidCreds = errors.New("invalid credentials")
var ErrInvalidSignature = errors.New("invalid signature")
var ErrInvalidToken = errors.New("invalid token")
var ErrExpiredToken = errors.New("expired token")
var ErrNotFound = errors.New("not found")
