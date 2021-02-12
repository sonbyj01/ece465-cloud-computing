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

<!-- TODO -->

---

### How the Distributed Algorithm Works (in more depth)

<!-- TODO: include system diagram -->

---

### Results

<!-- TODO -->

---

### Future Work
Gebremedhin et al. (2005) achieved almost linear speedups.
The goal for the second revision of this project is to provide further
optimization to achieve a better theoretical result. There are many different
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
process of distributing config files and start the clients and server. We also
need to write a lot of tests (similar to `proj1_tests.go`) to make sure that
everything is working as expected.


[gam2000]: http://www.ii.uib.no/~fredrikm/fredrik/papers/Concurrency2000.pdf
[gam2005]: https://cscapes.cs.purdue.edu/coloringpage/abstracts/euro05.pdf
