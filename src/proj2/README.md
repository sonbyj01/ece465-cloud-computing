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
Running `make` without a target **in the top-level directory of this repo** will
display a help menu:
```text
$ make
Usage: make [COMMAND], where COMMAND is one of the following:
        server: build server
        client: build client
        run-server: run server (build if necessary)
        run-client: run client (build if necessary)
        clean: clean built files
        logclear: clear logfiles
        refresh: runs targets clean logclear server client
```
The first two commands will build target files to [`target`](../../target).
Namely, it will build an executable to `target/server_{VERSION}` or
`target/client_{VERSION}`, as well as a symlink called `target/server_latest`
or `target/client_latest`. (Note that the versioned executable
file should be copied to a worker compute node rather than the symlink;
the symlink only exists for convenience.)

The `run-server` and `run-client` commands invoke `target/server_latest` or
`target/client_latest` with default parameters. Again, this is for convenience
and will likely not be the case -- if you need custom parameters, run the
built executables in the `target` directory.

##### Sample Run
To run multiple clients on the same device, you can use the
[`sample_run.sh`](../../sample_run.sh) located in the top-level directory.
Sample usage (from the top-level directory of this repo):
```bash
$ make refresh && ./sample_run res/sample10.graph
```

**Note for submission 2a**: The algorithm is still a little buggy and you may
have to run this a few times to get a clean output. If the algorithm does not
cleanly terminate (e.g., in the case of a short write), you may also want to run
```bash
$ pkill client_latest;pkill server_latest
```

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
Each line will contain the address of one worker. This file will be fed to the
server.
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

### Progress & Results

##### Project 2(a) -- Handshake complete, algorithm prone to race conditions
We were able to implement the handshake successfully, but encountered multiple
bugs when actually implementing the distributed algorithm. These bugs
are related to race conditions (e.g., negative WaitGroup counters) and
a short write condition.

Currently, the binary is very buggy and may take a few runs until it doesn't
cause panics due to these bugs. You can try to build using the sample
instructions listed above but it may take 4-5 tries to get a suitable output.
A sample output when it does work correctly, for comparison purposes, can be
found at [`res/sample10.log`](../../res/sample10.log). This colored the graph
as follows (see graph file format):
```text
10
0:2;1,9,7
1:1;0,8,4
2:1;6
3:1;4
4:2;1,3
5:1;8,7
6:0;2
7:3;8,0,5
8:0;1,5,7
9:1;0
```
It is easy to verify that this coloring is correct.

Our main goal for the second iteration is to get this working, and if possible,
to measure the speedup this achieves over sequential operation and single-node,
multithread operation (if applicable). Since the handshake works properly, this
will probably mean rewriting most of the algorithm with very thorough testing
to eliminate the race conditions and the short write condition.

---

### Project 2b
For checkpoint 2b, the first order of business was to fix the bugs in 2a, so
that the algorithm runs reliably. Unfortunately, that did not leave much time
for optimizations or deploying to AWS, which will be for future assignments.

A full write-up for 2b, as well as general documentation and visuals on the
distributed algorithm, can be found [here][2bdoc].


[gam2000]: http://www.ii.uib.no/~fredrikm/fredrik/papers/Concurrency2000.pdf
[gam2005]: https://cscapes.cs.purdue.edu/coloringpage/abstracts/euro05.pdf
[2bdoc]: http://files.lambdalambda.ninja/reports/20-21_spring/ece465_proj2b.lam_son.pdf