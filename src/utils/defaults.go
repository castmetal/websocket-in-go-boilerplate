package utils

type params struct {
	PORT        string
	ENV         string
	AUTH_HEADER string
}

var defaultParams = params{
	PORT:        "8000",
	ENV:         "staging",
	AUTH_HEADER: "X-Auth",
}

var portEnvVar = GetEnvVariable("PORT")
var port = *CoalesceString(&portEnvVar, &defaultParams.PORT)

var envVar = GetEnvVariable("ENV")
var env = *CoalesceString(&envVar, &defaultParams.ENV)

var authHeaderVar = GetEnvVariable("AUTH_HEADER")
var authHeader = *CoalesceString(&authHeaderVar, &defaultParams.AUTH_HEADER)

var SystemParams = params{
	PORT:        port,
	ENV:         env,
	AUTH_HEADER: authHeader,
}
