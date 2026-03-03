##
## D&GINE Project, 2026
## API
## File description:
## Makefile
##

GO	=	go
GOFLAGS = -ldflags="-s -w"

BINARY_NAME	=	api

## NORMAL MAIN COMPILATION
all: $(BINARY_NAME)

$(BINARY_NAME):	$(SRC) build run

build: $(SRC)
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) ./src/main.go

run:
	./$(BINARY_NAME)

re:     fclean all

## CLEANING RULES
clean:
	$(GO) clean

fclean: clean
	$(RM) $(BINARY_NAME)

.PHONY: all build run re clean fclean