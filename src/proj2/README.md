# ECE465 Project 2: Multi-Node, Multi-Thread Graph Coloring
### Jonathan Lam & Henry Son

---

### Algorithm Overview
For a general overview of graph coloring, see the
[Project 1 README](../proj1/README.md). The same general algorithm from
[Gebremedhin and Manne (2000)][gam2000] from Project 1 is used on each node, and the
distributed algorithm from [Gebremedhin et al. (2005)][gam2005] is used to coordinate
distributed nodes.

---

### Build Instructions
The top-level [Makefile](../../Makefile) has several instructions to build
the executables and run defaults. Running `make` (without a target) in the
terminal will provide some aid:
```text
$ make
Usage: make [COMMAND], where COMMAND is one of the following:
        server: build server
        client: build client
        run-server: run server (build if necessary)
        run-client: run client (build if necessary)
        clean: clean built files
```
The first two commands will build target files to [`/target`](../../target).
Namely, it will build an executable to `target/server_{VERSION}` or
`target/client_{VERSION}`, as well as a symlink called `target/server_latest`
or `target/client_latest`. (Note that the versioned executable
file should be copied to a worker compute node rather than the symlink;
the symlink only exists for convenience.)

The `run-server` and `run-client` commands invoke `target/server_latest` or
`target/client_latest` with default parameters. Again, this is for convenience
and will likely not be the case -- if you need custom parameters, run the
built executables in the `/target` directory.

---

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

---

### Distributed Environment Setup
To run and test the first revision, we run all of the nodes on the same host
under different ports. We simply did not have enough time to test this across
true multi-node environments, which is something we aim to do for the second
revision.

---

### How the Distributed Algorithm Works (in more depth)

##### Initial Server/Worker Handshake
The server reads in the configuration `*.nodes` file indicating where each
node is listening (IP address and port of each worker node). The server
attempts to make a handshake with each worker node, after which it sends
information about the other nodes so that they can each create peer-to-peer
nodes amongst themselves (the graph of worker nodes forms a complete graph).

(Note about revision 1: We were not able to implement the Server-Worker
handshake in the manner described above. A simplified version, where each
worker is given the `*.nodes` configuration file, is used instead for this
revision.)

##### Distributing the Graph
The graph is read in from a file on the server, and partitioned into
subgraphs comprising vertices with adjacent indices. (This is the simplest
graph partitioning scheme, albeit not nearly optimal -- see the comments in
the Future Work section.) These subgraphs are sent off to their designated
worker nodes, and then the worker nodes are told to begin the algorithm.

##### Worker Thread: Single-Node Multi-Thread
This proceeds in much the same way as in Project 1. The difference is that
each node is assigned a subgraph, which has many edges that lead to vertices
in other subgraphs. To handle this, every time vertex on the edge of a node's
subgraph is (re)colored, its neighbors are notified of its new color. (This
process is buffered for performance, and it is arbitrarily chosen that only
the higher-numbered vertex is notified of coloring changes so that only one
of the nodes gets recolored.) Each node has a goroutine listening for
notifications from other threads while it is performing the coloring.

After a node finishes coloring, it notifies the other nodes and waits for all
other nodes to finish coloring. This forces all nodes to perform each step
in lockstep; this is the synchronized algorithm proposed in Gebremedhin et al.
(2005), which we adopt since it is simpler than the asynchronous version.

When all nodes are finished coloring for the current coloring iteration, then
conflict detection occurs in the same way as in Project 1. Note that this
requires no inter-node communication, since all inter-node communication happens
during the speculative coloring stage.

##### Algorithm Completion
After a node finishes coloring, it can begin the next iteration (without
requiring other nodes to finish). When there are zero conflicts (no nodes to
be recolored), then it broadcasts a different signal to all of the other nodes.
When all workers are done, then the algorithm terminates and the workers stream
their subgraphs back to the server for reconstitution and verification of
proper coloring.

<!-- TODO: include system diagram -->

---

### Results

<!-- TODO -->

---

### Future Work
Gebremedhin et al. (2005) achieved almost linear speedups.
The goal for the second revision of this project is to provide further
optimization to achieve a better theoretical result. There are many
directions we can explore for the second
revision, e.g., switching away from the most trivial first-fit (FF) color-choice
algorithm, adjusting buffering sizes, optimizing datatype sizes.

What would most
likely cause the largest difference in results is to choose a real-world graph
that is well-partitioned (i.e., so that a good partitioning greatly reduces
the number of inter-node edges from the uniform case), as well as the
partitioning. Gebremedhin et al. (2005) use a real-world graph and the METIS
graph partitioning tool. Since graph partitioning (to minimize cross-edges)
is a NP-hard problem, we would also likely use a tool like this; however,
there is no implementation of METIS for Golang, so we would have to use another
tool (or pre-process the graph before use in Golang).

Other optimizations may be on the network stack, since the networking is
relatively primitive. We use raw TCP sockets and channels with `gob`-ed data
(i.e., "serialized" similar to Python's pickle), and send data over slices --
it may be faster to send simpler fixed-size packets of a reasonable buffer size
to improve performance.

We also need to continue working on our tooling. For now, we have a very basic
outline in [`automate_start.sh`](../../automate_start.sh) to automate the
process of distributing config files and start the clients and server. There
are also a lot of missing utilities simply due to the fact that we didn't finish
everything on time, such as reconstituting and checking the graph on the server
on completion, etc.

Also, some general housekeeping items that we ran out of time to do for the
first revision:
- Unit tests
- Code cleanup and commenting to please `golint src/**`


[gam2000]: http://www.ii.uib.no/~fredrikm/fredrik/papers/Concurrency2000.pdf
[gam2005]: https://cscapes.cs.purdue.edu/coloringpage/abstracts/euro05.pdf
