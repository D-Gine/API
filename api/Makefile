##
## D&GINE Project, 2026
## API
## File description:
## Makefile
##

GO	=	go
GOFLAGS = -ldflags="-s -w"

BINARY_NAME	=	api

MAIN_FILE = ./src/main.go
DOCS_FOLDER = ./src/docs/

# SWAG = $(GO) run github.com/swaggo/swag/cmd/swag@latest
SWAG = swag

all: build

build:
	$(SWAG) init -g $(MAIN_FILE) -o $(DOCS_FOLDER)
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(MAIN_FILE)

run: build
	./$(BINARY_NAME)

re:     fclean all

## CLEANING RULES
clean:
	$(GO) clean

fclean: clean
	$(RM) $(BINARY_NAME)

.PHONY: all build run re clean fclean