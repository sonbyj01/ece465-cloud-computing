build:
	go build ./src/proj2/server
	go build ./src/proj2/client
run-server:
	./server -config config.nodes
run-client:
	./client -intf ens33 -port 8007
