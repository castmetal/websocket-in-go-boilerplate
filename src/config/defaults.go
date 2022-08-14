package config

const (
	PORT             = "8000"
	ENV              = "staging"
	AUTH_HEADER      = "X-Auth"
	SERVER_TYPE      = "ws"
	SALT_CHARSET     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	DB_HOST          = "localhost"
	DB_PORT          = "5432"
	DB_USER          = "myusername"
	DB_PASSWORD      = "mypassword"
	DB_TIME_ZONE     = "America/Sao_Paulo"
	DB_DATABASE_NAME = "public"
)

type Config struct {
	PORT             string
	ENV              string
	AUTH_HEADER      string
	SERVER_TYPE      string
	SALT_CHARSET     string
	DB_HOST          string
	DB_PORT          string
	DB_USER          string
	DB_PASSWORD      string
	DB_TIME_ZONE     string
	DB_DATABASE_NAME string
}

var cfg = Config{
	PORT,
	ENV,
	AUTH_HEADER,
	SERVER_TYPE,
	SALT_CHARSET,
	DB_HOST,
	DB_PORT,
	DB_USER,
	DB_PASSWORD,
	DB_TIME_ZONE,
	DB_DATABASE_NAME,
}

var portEnvVar = GetEnvVariable("PORT")
var envVar = GetEnvVariable("ENV")
var authHeaderEnvVar = GetEnvVariable("AUTH_HEADER")
var serverTypeEnvVar = GetEnvVariable("SERVER_TYPE")
var saltCharsetEnvVar = GetEnvVariable("SERVER_TYPE")
var dbHostEnvVar = GetEnvVariable("DB_HOST")
var dbUserEnvVar = GetEnvVariable("DB_USER")
var dbPortEnvVar = GetEnvVariable("DB_PORT")
var dbPasswordEnvVar = GetEnvVariable("DB_PASSWORD")
var dbTimeZoneEnvVar = GetEnvVariable("DB_TIME_ZONE")
var dbDatabaseEnvVar = GetEnvVariable("DB_DATABASE_NAME")

var SystemParams = Config{
	PORT:             *CoalesceString(&portEnvVar, &cfg.PORT),
	ENV:              *CoalesceString(&envVar, &cfg.ENV),
	AUTH_HEADER:      *CoalesceString(&authHeaderEnvVar, &cfg.AUTH_HEADER),
	SERVER_TYPE:      *CoalesceString(&serverTypeEnvVar, &cfg.SERVER_TYPE),
	SALT_CHARSET:     *CoalesceString(&saltCharsetEnvVar, &cfg.SALT_CHARSET),
	DB_HOST:          *CoalesceString(&dbHostEnvVar, &cfg.DB_HOST),
	DB_USER:          *CoalesceString(&dbUserEnvVar, &cfg.DB_USER),
	DB_PORT:          *CoalesceString(&dbPortEnvVar, &cfg.DB_PORT),
	DB_PASSWORD:      *CoalesceString(&dbPasswordEnvVar, &cfg.DB_PASSWORD),
	DB_TIME_ZONE:     *CoalesceString(&dbTimeZoneEnvVar, &cfg.DB_TIME_ZONE),
	DB_DATABASE_NAME: *CoalesceString(&dbDatabaseEnvVar, &cfg.DB_DATABASE_NAME),
}
