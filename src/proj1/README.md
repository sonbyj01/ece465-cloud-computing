# ECE465 Project 1: Graph Coloring

---

### Project Overview
Given an arbitrary undirected (finite) graph, we can "color" each node in
a way such that neighbors do not have the same color. The **graph coloring
problem** involves finding such a coloring of a graph, or the minimum
number of colors needed to fulfill this criteria (known as the "chromatic
number" of the graph). This is closely related to the independent sets
problem (each set of nodes sharing the same color is an independent set),
and has applications in distributed resource allocation (e.g., for
compiler or distributed algorithm optimization). There are also many special
cases of the problem, such as in the case of two-coloring (bipartite graphs),
ring graphs (which can be proved to have a chromatic number of 3), and
most famously planar graphs (which have been proved to have a chromatic number
of 4, such as in the case of maps). It can also be used to solve Sudoku,
whose constraints can be phrased as a graph coloring problem.

This is a brief overview of the problem; a more complete and very well-written
overview of the graph coloring problem, its applications, known algorithms
(sequential/parallel (shared-memory and message-passing based)), and analysis
of the algorithms can be found in [Gebremedhin's thesis][gebre].

---

### Build Instructions

See [the parent README](../../README.md).

---

### Algorithm Overview
While finding the (true) chromatic number is an
NP-complete problem (and thus the only exact algorithms are
exponential-time algorithms), there exist efficient randomized, distributed
algorithms that can approximate the chromatic number (and perform fairly
well in real-world tasks). There are a number of known sequential and parallel
algorithms

The algorithm we choose is called the
Gebremedhin-Manne algorithm, which is described in [this OSTI presentation][1];
this takes a greedy approach by speculatively choosing the lowest
whole number (color) not taken by any of its neighbors, for multiple
nodes in parallel. Since this step is done without any data synchronization,
it is fast but may result in some inconsistencies (conflicts). Thus
after performing this step on every node, a synchronization step
(a [barrier][wg]) is performed to detect improper colorings (data
inconsistencies due to race conditions), and then we repeat the same procedure
again only on the nodes which had conflicts until no conflicts arise.

The [paper describing this method][5] goes into more detail. We did not
implement this algorithm in its entirety, but rather used the simplified version
presented in the OSTI presentation. We may implement more of this algorithm
in its entirety for revision 2 of this project.

##### Notes on Implementation and Tests
A sequential and parallel version were written for this project. They both
employ the same general idea, except that the parallel version runs the
speculative coloring stage and conflict detection stages in parallel, and
repeatedly colors the conflicted nodes until all nodes are colored (conflicts
are not possible in the sequential case).

For the parallel stage, we spawned a new goroutine for each node to be
colored/conflict-checked. While this might be inefficient using threading
libraries in most languages, [goroutines are more lightweight than OS threads
and are designed to be spawned en masse][goroutines]. However, it would be a
good idea to experiment with spawning fewer goroutines (e.g., up to the number
of logical processors on the compute node) to see if this affects the
performance.

For testing, we created functions in `graph/generate_graphs.go` that can
generate complete graphs, ring graphs, and graphs with a given branching factor.
The former two are special cases in the graph coloring problem with known
solutions (complete graphs have a chromatic number of N-1 and ring graphs have
a chromatic number of either 2 or 3).
We use the latter for our simulations because it is more generic than the
other two; however, it generates uniformly-distributed graphs, which are unlike
many real-world graphs with nonuniform distributions, e.g., small-world networks
and topologies with clustering patterns. The generation of the random graph
is also multithreaded to reduce graph generation time. (Graph generation
actually takes longer than graph coloring in our tests because of the number
of potential edges it loops through, so this saves a lot of time in our
benchmarks.)

Our test and benchmark framework is based on Golang's [`testing`][testing]
package. A new graph is generated before each benchmark, and the benchmarking
timer is paused during graph generation. See [the parent
README](../../README.md) for details on how to run the benchmarks.

##### Alternative Parallel Algorithms
(A more complete overview of graph coloring algorithms can be found in
[Gebremedhin's thesis][gebre].)

The reason we chose this algorithm was that it is simple and empirically fast;
it is the choice of algorithm in the OSTI presentation. All we do is
speculatively and greedily color the nodes in parallel and cross our fingers
that no conflicts occur (without requiring synchronization), and then check for
conflicts to resolve. Since the chance of conflict is very low, we expect that
the conflict resolution costs less than the overhead of synchronization. A
variation on this is to perform locking when coloring any vertex.

Another approach, suggested by [Jones and Plassman][jp], is based on the fact
that the graph coloring problem is closely related to the problem of
[independent sets][is]: in particular, each set of nodes that share the same
color (in a valid coloring) forms an independent set. Finding the largest
independent set in a graph (the maximum independent set problem) is an NP-hard
problem, but there are polynomial-time distributed algorithms to find large
independent sets, e.g., [Luby's algorithm][luby] and [an algorithm by
Schneider and Wattenhofer][3]. A basic distributed algorithm suggests coloring
an entire maximal independent set the same color, but the algorithm by Jones and
Plassman use an independent set to choose the order with which to color nodes.

Another approach is taken using [Kuhn-Wattenhofer color reduction][2], which
performs a divide-and-conquer approach, not by partitioning the vertices, but
by partitioning the color classes. We didn't find a comparison between this
method and the previous ones mentioned, so we are not sure how it performs.
However, a cursory look at the pseudocode in this article indicates that it
doesn't attempt to find a good approximation of the chromatic number (only
reducing the number of colors to something less than the worst-case bound), but
that  it rather focuses on speed. The analysis provided by the Jones-Plassman
and Gebremedhin algorithms in their respective papers indicate that they provide
a fairly good approximation of the chromatic number.

##### Runtime Analysis and Experimental Results

TODO: need to write more tests to explore this

##### Next Steps: Scaling Up to Multi-Node
(For project 2)

The variation of Gebremedhin's algorithm that we used, largely based on the
pseudocode from OSTI, has the downfall that it is relies on a
shared-memory architecture. The problem is that we lose a centralized,
high-speed memory when switching to a multi-node architecture. In a multi-node
environment, we have to rely on message-passing between nodes (or to and from
a distributed storage node), which would incur much higher latency and probably
bottleneck our algorithm.

It would be preferable to limit inter-node communication since it is so much
slower than accessing a local compute node's RAM. In general, one way to do this
is to [partition the graph][partition] in such a way that the number of edges
connecting the subgraphs are minimized; but this is an NP-hard problem, so we
also look for an approximation to this. Jones and Gebremedhin suggest performing
a simple partitioning (without any attempt to minimize the number of connecting
edges) and performing a two-stage coloring:

1. Color the nodes that are adjacent to nodes on another graph (this may either
   be parallelized but require a significant degree of message passing, or it
   may be simply performed on a single node)
2. Color the nodes within each subgraph (this can be entirely parallelized)

where each subgraph in step 2 may be colored with any single-node, multithreaded
coloring algorithm. Thus this general idea can be implemented in many different
particular ways depending on the choice of sub-algorithm.

Another consideration is that we have not explored this algorithm in the space
of very-large graphs, especially graphs that will no longer fit into main
memory. For example, in these tests, a rough estimate (from looking at `top`'s
output) is that a `graph.Graph` instance generated with `graph.NewRandomGraph`
with 50000 nodes and an average branching factor of 250 takes over 4Gb of RAM;
if we were to deal with much larger graphs, it would be problematic to try
to fit them into main memory (assuming that a single compute node is limited to
4 to 16GB of RAM). This is something we have to look into for future revisions
of our algorithm.

---

### Extra Resources
- [A smaller exponential time bound for exact graph coloring (using
  small maximal independent sets)][4]
- [Graph coloring on the GPU][6]: has a good literature review of distributed
  graph coloring algorithms

[1]: https://www.osti.gov/servlets/purl/1246285
[2]: https://stanford.edu/~rezab/classes/cme323/S16/projects_reports/bae.pdf
[3]: https://tik-old.ee.ethz.ch/file//2be1291694b1730bba83f7fa18d9e0f2/podc08SW.pdf
[4]: https://arxiv.org/pdf/cs/0011009.pdf
[5]: https://cscapes.cs.purdue.edu/coloringpage/abstracts/euro05.pdf
[6]: https://people.eecs.berkeley.edu/~aydin/coloring.pdf 
[wg]: https://stackoverflow.com/a/22697521/2397327
[gebre]: https://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.126.882&rep=rep1&type=pdf
[partition]: https://en.wikipedia.org/wiki/Graph_partition
[jp]: https://pdfs.semanticscholar.org/a2a5/43255d7db24dbd71798d87eac5cf30883b3a.pdf
[is]: https://en.wikipedia.org/wiki/Independent_set_(graph_theory)
[luby]: https://en.wikipedia.org/wiki/Maximal_independent_set#Random-selection_parallel_algorithm_%5BLuby%27s_Algorithm%5D
[goroutines]: https://rcoh.me/posts/why-you-can-have-a-million-go-routines-but-only-1000-java-threads/
[testing]: https://golang.org/pkg/testing/