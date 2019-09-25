package common

type contextKey struct {
	name string
}

// UsernameCtxKey defines username request context key.
var UsernameCtxKey = &contextKey{"Username"}

// ShareCtxKey defines share request context key.
var ShareCtxKey = &contextKey{"Share"}
