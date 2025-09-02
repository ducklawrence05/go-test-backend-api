package jwtpurpose

type JWTPurpose string

const (
	Access   JWTPurpose = "access"
	Refresh  JWTPurpose = "refresh"
	Register JWTPurpose = "register"
	Restore  JWTPurpose = "restore"
)
