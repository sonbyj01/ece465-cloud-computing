# Makefile for Project 2 (distributed coloring algorithm)
# @author Jonathan Lam <jlam55555@gmail.com>
# @author Henry Son <sonbyj01@gmail.com>

OUTDIR=./target

# set gopath
GOENV:=GOPATH=$(CURDIR)

# versioning is created with git version and date
GIT_BRANCH:=$(shell git status|head -n 1|awk '{print $$3}')
GIT_REV:=$(shell git log|head -n 1|awk '{print $$2}')
TIMESTAMP:=$(shell TZ=UTC date "+%Y%m%d-%H%M%S")
VERSION:=$(TIMESTAMP)_$(GIT_BRANCH)_$(GIT_REV)

# server vars
SERVER_BIN=server
SERVER_FLAGS=-config config.nodes

# client vars
CLIENT_BIN=client
CLIENT_FLAGS=-port 8007

.PHONY:
default:
	@# @ hides command: https://stackoverflow.com/a/9967125/2397327
	@echo "Usage: make [COMMAND], where COMMAND is one of the following:"
	@echo "	run-server: build and run the server"
	@echo "	run-client: build and run the client"
	@echo "	clean: clean built files"

.PHONY:
run-server:
	$(GOENV) go build -o $(OUTDIR)/$(SERVER_BIN)_latest ./src/proj2/server
	cp $(OUTDIR)/$(SERVER_BIN)_latest $(OUTDIR)/$(SERVER_BIN)_$(VERSION)
	$(OUTDIR)/$(SERVER_BIN)_latest $(SERVER_FLAGS)

.PHONY:
run-client:
	$(GOENV) go build -o $(OUTDIR)/$(CLIENT_BIN)_latest ./src/proj2/client
	cp $(OUTDIR)/$(CLIENT_BIN)_latest $(OUTDIR)/$(CLIENT_BIN)_$(VERSION)
	$(OUTDIR)/$(CLIENT_BIN)_latest $(CLIENT_FLAGS)

.PHONY:
clean:
	rm -rf pkg target
