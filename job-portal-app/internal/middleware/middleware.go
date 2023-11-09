package middleware

import (
	"job-portal-api/internal/auth"
)

type Mid struct {
	auth auth.Authentication
}

func NewMiddleware(a auth.Authentication) (Mid, error) {
	return Mid{
		auth: a,
	}, nil
}
