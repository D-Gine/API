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
GIN_MODE=debug
TOKEN_KEY=[TOKEN_KEY : string]
ALLOWED_ORIGINS=http://localhost:8081,http://127.0.0.1:8081
TOKEN_DOMAIN=
TRUSTED_PROXIES=127.0.0.1

# NETWORKS
BACK_NETWORK=dengine-back-net

# VOLUMES
API_VOLUME=dengine-api-volume

# FILESERVER
IMAGES_PATH=/img
USERS_IMAGES_PATH=/img/users/
```

### Launch
#### Development
```
docker compose up --build -d
```

## Doc
You have access to routes documentations via Swagger, it must be and shall be updated with each new feature on the API
<br>To access the swagger, go to [THIS LINK](http://localhost:8080/swagger/index.html) once you started the API on your computer

## Participate
You can add features whenever you feel like, but keep in mind to document them (in the swagger if they are routes)

## Useful knowledge for devs
### Swagger
init path
```
export PATH=$(go env GOPATH)/bin:$PATH
```

set docs files
```
swag init -g <main.go file> -o <doc folder dest>
```
