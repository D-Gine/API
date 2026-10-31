GO	=	go
GOFLAGS = -ldflags="-s -w"

BINARY_NAME	=	api

MAIN_FILE = ./src/main.go
DOCS_FOLDER = ./src/docs/

SWAG = swag
SWAG_CMD = $(SWAG) init -g $(MAIN_FILE) -o $(DOCS_FOLDER)

RUN_MODE ?= AUTO
RUN_MODE_UP = $(shell echo '$(RUN_MODE)' | tr '[:lower:]' '[:upper:]')

all: build

swagger:
	$(SWAG_CMD)

build: swagger
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(MAIN_FILE)

run:
ifeq ($(RUN_MODE_UP),BIN)
	@$(MAKE) --no-print-directory run-bin
else ifeq ($(RUN_MODE_UP),AUTO)
	@$(MAKE) --no-print-directory run-auto
else
	$(error invalid RUN_MODE given : '$(RUN_MODE)'; should be 'AUTO' or 'BIN')
endif

run-bin: build
	./$(BINARY_NAME)

run-auto:
	reflex -s -r '\.go$$' -R '^src/docs/' -- \
	sh -c '$(SWAG_CMD) && exec $(GO) run $(MAIN_FILE)'

re:     fclean all

clean:
	$(GO) clean

fclean: clean
	$(RM) $(BINARY_NAME)

.PHONY: all swagger build run run-bin run-auto re clean fclean