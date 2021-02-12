TODO

### Build and Run Instructions
```bash
$ go build ./src/proj2/server
$ go build ./src/proj2/client
```

```bash
$ ./server --config config.nodes
$ ./client --intf enps025 --port 5000
```

### File Formats

##### Graph Description Files (*.graph)
First line is number of vertices. Each line after that indicates one vertex:
its value and adjacency list (comma-separated).

Sample file for a graph with 3 nodes, with values 1, 2, and -2:
```text
3
0;1,2
2;0
-2;0
```

##### Node Configuration Files (*.node)
Each line will contain the address of one slave. This file will be fed to the
master. (Each slave will also be fed the master's IP address.)
```text
10.0.0.53:5000
10.0.0.54:5000
10.0.0.55:5000
10.0.0.56:5000
```