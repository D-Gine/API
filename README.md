# D&Gine API

## Prerequisites
To run the API using the docker-compose, you need to have an existing external network with the database connected
<br>i.e. launched [THIS PROJECT](https://github.com/D-Gine/Database)

## Build and Run with
### Env
For this project you need the following environment variables:
```
# DATABASE
DB_HOST=dengine-db
DB_PORT=5432
DB_USER=[DB_USER : string]
DB_PWD=[DB_PWD : string]
DB_NAME=dengine

# API
VERSION=1.3.0
API_PORT=[API_PORT : int]
API_SRC=./api/src
RUN_MODE=AUTO
GIN_MODE=debug
TOKEN_KEY=[TOKEN_KEY : string]
ALLOWED_ORIGINS=http://localhost:8081,http://127.0.0.1:8081
TOKEN_DOMAIN=
TRUSTED_PROXIES=127.0.0.1

# NETWORKS
BACK_NETWORK=dengine-back-net
FRONT_NETWORK=dengine-front-net

# VOLUMES
API_VOLUME=dengine-api-volume

# FILESERVER
IMAGES_PATH=/img
USERS_IMAGES_PATH=/img/users/
```

**RUN_MODE**: Determines if the API will be run as **binary** (for prod) or in **REFLEX** (for dev). Reflex provides automatic reload of the files onto the container. 

### Launch
#### Development
You need to have the **environment variables** set in `.env` file at the same place as the `compose.yaml`
```
docker compose up --build -d
```

## Doc
You have access to routes documentations via Swagger, it must be and shall be updated with each new feature on the API
<br>To access the swagger, go to [THIS LINK](http://localhost:8090/scalar) once you started the API on your computer

## Participate
You can add features whenever you feel like, but keep in mind to document them (in the swagger if they are routes)

## Useful knowledge for devs
### Unit tests

Unit tests in golang are automatically recognized and launched with the following command, executed where the `main.go` file sits
```sh
go test ./...
```
Basic golang **testing** module working :
- It browses automatically files ending with `_test.go` (ex : `api/src/controllers/auth/auth_test.go`)
- Inside the files, it browses automatically functions starting with `Test` (ex : `Test_CreateToken`)

A few rules for this project specifically :
1. Each package needs to have one unit tests file with the syntax `<package>_test.go` (ex: `auth_test.go`)
2. Each function tested must have only one starting point function of testing called `Test_<Function>`
3. If you need multiple tests cases for a function, you will use subfunctions. subfonctions format is `Sub_<Function>_<CaseTested>`

Example of subfunctions to test `ExtractFromToken`:
```go
func Test_ExtractFromToken(t *testing.T) { // tests starting point function
	t.Run("extract from token", Sub_ExtractFromToken) // call to subfunction to test case OK
	t.Run("with bearer prefix", Sub_ExtractFromToken_WithBearerPrefix) // call to subfunction to test case OK with bearer
	t.Run("invalid token", Sub_ExtractFromToken_Invalid) // call to subfunction to test case of invalid token
	t.Run("missing or invalid claim", Sub_ExtractFromToken_MissingOrInvalidClaim) // call to subfunction to test case of missing or invalid claim
}
```

---

### Swagger
Swagger is now automatically reloaded on every build of the API

init path
```
export PATH=$(go env GOPATH)/bin:$PATH
```

set docs files
```
swag init -g <main.go file> -o <doc folder dest>
```
