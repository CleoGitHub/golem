package httpclient

import "errors"

var (
	ErrUnexpectedStatus  = errors.New("unexpected response status")
	ErrTokenEmpty        = errors.New("token empty in response")
	ErrRefreshTokenEmpty = errors.New("refresh token empty in response")
)
