# Makefile for Project 2 (distributed coloring algorithm)
# @author Jonathan Lam <jlam55555@gmail.com>
# @author Henry Son <sonbyj01@gmail.com>

### CONFIGURABLE VARIABLES

# server vars
SERVER_BIN=server
SERVER_FLAGS=-config config.nodes

# client vars
CLIENT_BIN=client
CLIENT_FLAGS=-port 8007 -config config.nodes

# build directory
OUTDIR=./target

### DON'T MODIFY ANYTHING BELOW THIS POINT

# set gopath
GOENV:=GOPATH=$(CURDIR)

# versioning is created with git version and date
GIT_BRANCH:=$(shell git status|head -n 1|awk '{print $$3}')
GIT_REV:=$(shell git log|head -n 1|awk '{print $$2}')
TIMESTAMP:=$(shell TZ=UTC date "+%Y%m%d-%H%M%S")
VERSION:=$(TIMESTAMP)_$(GIT_BRANCH)_$(GIT_REV)

# filenames
SERVER_VERFILE:=$(OUTDIR)/$(SERVER_BIN)_$(VERSION)
SERVER_LATEST:=$(OUTDIR)/$(SERVER_BIN)_latest
CLIENT_VERFILE:=$(OUTDIR)/$(CLIENT_BIN)_$(VERSION)
CLIENT_LATEST:=$(OUTDIR)/$(CLIENT_BIN)_latest

# Note: @ hides command: https://stackoverflow.com/a/9967125/2397327
.PHONY: default
default:
	@echo "Usage: make [COMMAND], where COMMAND is one of the following:"
	@echo "	server: build server"
	@echo "	client: build client"
	@echo "	run-server: run server (build if necessary)"
	@echo "	run-client: run client (build if necessary)"
	@echo "	clean: clean built files"

# canned recipes because we may want to always rebuild server or not
define BUILD_SERVER=
$(GOENV) go build -o $(SERVER_VERFILE) ./src/proj2/server
ln -sf $(CURDIR)/$(SERVER_VERFILE) $(SERVER_LATEST)
endef
define BUILD_CLIENT=
$(GOENV) go build -o $(CLIENT_VERFILE) ./src/proj2/client
ln -sf $(CURDIR)/$(CLIENT_VERFILE) $(CLIENT_LATEST)
endef

# build server (always rebuild since no easy way to check all deps)
.PHONY: server
server:
	$(BUILD_SERVER)

# non-phony target to rebuild server when necessary
$(SERVER_LATEST):
	$(BUILD_SERVER)

# run server (only build if necessary)
run-server: $(SERVER_LATEST)
	$(SERVER_LATEST) $(SERVER_FLAGS)

# build client (always rebuild)
.PHONY: client
client:
	$(BUILD_CLIENT)

# target to only rebuild client if necessary
$(CLIENT_LATEST):
	$(BUILD_CLIENT)

# run client
run-client: $(CLIENT_LATEST)
	$(CLIENT_LATEST) $(CLIENT_FLAGS)

# remove built executables
.PHONY: clean
clean:
	rm -rf pkg target
