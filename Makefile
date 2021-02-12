run-server:
	go build ./src/proj2/server
	./server -config config.nodes

run-client:
	go build ./src/proj2/client
	./client -intf ens33 -port 8007

clean:
	rm -rf server client
