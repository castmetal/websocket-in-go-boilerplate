package config

const (
	PORT        = "8000"
	ENV         = "staging"
	AUTH_HEADER = "X-Auth"
	SERVER_TYPE = "ws"
)

type Config struct {
	PORT        string
	ENV         string
	AUTH_HEADER string
	SERVER_TYPE string
}

var cfg = Config{
	PORT,
	ENV,
	AUTH_HEADER,
	SERVER_TYPE,
}

var portEnvVar = GetEnvVariable("PORT")
var envVar = GetEnvVariable("ENV")
var authHeaderEnvVar = GetEnvVariable("AUTH_HEADER")
var serverTypeEnvVar = GetEnvVariable("SERVER_TYPE")

var SystemParams = Config{
	PORT:        *CoalesceString(&portEnvVar, &cfg.PORT),
	ENV:         *CoalesceString(&envVar, &cfg.ENV),
	AUTH_HEADER: *CoalesceString(&authHeaderEnvVar, &cfg.AUTH_HEADER),
	SERVER_TYPE: *CoalesceString(&serverTypeEnvVar, &cfg.SERVER_TYPE),
}