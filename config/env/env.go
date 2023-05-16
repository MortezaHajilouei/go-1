package env

type EnvMode string

const (
	PRODUCTION  EnvMode = "prod"
	STAGE       EnvMode = "stage"
	DEVELOPMENT EnvMode = "dev"
	TEST        EnvMode = "test"
)
