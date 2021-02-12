OUTDIR=./target

# server vars
SERVER_BIN=server
SERVER_FLAGS=-config config.nodes

# client vars
CLIENT_BIN=client
CLIENT_FLAGS=-port 8007

.PHONY:
run-server:
	go build -o $(OUTDIR)/$(SERVER_BIN) ./src/proj2/server
	$(OUTDIR)/$(SERVER_BIN) $(SERVER_FLAGS)

.PHONY:
run-client:
	go build -o $(OUTDIR)/$(CLIENT_BIN) ./src/proj2/client
	$(OUTDIR)/$(CLIENT_BIN) $(CLIENT_FLAGS)

.PHONY:
clean:
	rm -rf pkg target
