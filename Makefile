run-server:
	go build ./src/proj2/server
	./server -config config.nodes

run-client1:
	go build ./src/proj2/client
	./client -port 8007 -config config.nodes

run-client2:
	./client -port 8008

clean:
	rm -rf server client
