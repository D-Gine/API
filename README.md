# D&Gine API

## Build and Run with

For this project you need the following environment variables:
```
# DATABASE
DB_HOST=dengine-db
DB_PORT=5432
DB_USER=[USER : string]
DB_PWD=[PWD : string]
DB_NAME=dengine

# API
API_PORT=[API_PORT : int]
GIN_MODE=debug
TOKEN_KEY=[TOKEN_KEY : string]
ALLOWED_ORIGINS=http://localhost:8081,http://127.0.0.1:8081
TOKEN_DOMAIN=
TRUSTED_PROXIES=127.0.0.1

# NETWORKS
BACK_NETWORK=dengine-back-net
```
