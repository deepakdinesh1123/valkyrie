package config

type StringContextKey string

var (
	UserKey StringContextKey = "user"
	AuthKey StringContextKey = "noauth"
)
