package jwtpurpose

type JWTPurpose string

const (
	JWTAccess   JWTPurpose = "access"
	JWTRefresh  JWTPurpose = "refresh"
	JWTRegister JWTPurpose = "register"
	JWTRestore  JWTPurpose = "restore"
)
