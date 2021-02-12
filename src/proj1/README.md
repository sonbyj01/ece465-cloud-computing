# ECE465 Project 1: Graph Coloring
### Jonathan Lam & Henry Son

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
algorithms, and we chose a fairly simple one for the first revision of this
project.

The algorithm we choose is called the
Gebremedhin-Manne algorithm, which was presented in [(Gebremedhin 2000)][gam]
and described more simply in [this OSTI presentation][1];
this takes a greedy approach by speculatively choosing the lowest
whole number (color) not taken by any of its neighbors, for multiple
nodes in parallel. Since this step is done without any data synchronization,
it is fast but may result in some inconsistencies (conflicts). Thus
after performing this step on every node, a synchronization step
(a [barrier][wg]) is performed to detect improper colorings (data
inconsistencies due to race conditions), and then we repeat the same procedure
again only on the nodes which had conflicts until no conflicts arise.

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
generate complete graphs, ring graphs, and graphs with a given average degree.
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
conflicts to resolve. Since the chance of conflict is very low, and true
synchronization would involve `O(Vd)=O(E)` mutex locks per iteration (which
we would expect to incur a high cost), we expect that
the conflict resolution costs less than the overhead of synchronization.

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
Let `V` be the number of nodes, and `d` be the average node degree (thus the
number of edges `E` is `Vd/2`).

The sequential algorithm loops over each vertex and finds the lowest valid color
for that node (i.e., the lowest color that has not been taken by any of
its neighbors). This is a linear search w.r.t. `b` for each node, and thus
the algorithm has time complexity `O(Vd)=O(E)`. The space complexity of the
algorithm is `O(1)`; we only allocate a (fixed size) buffer to keep track of
the colors of a node's neighbors.

As a very basic analysis of the parallel algorithm, if we assume:
- zero latency in thread scheduling
- zero latency in context switches
- all logical cores are in use during concurrent stages of the algorithm (e.g.,
  all threads in a pool finish at the same time)

then the first stage (speculative coloring) is an ideally-parallelized version
of the sequential algorithm, and takes `O(Vd/N)` time, if we let `N` be the
number of parallel threads that can be run at a time (i.e., the number of
logical cores).

The second stage takes the same amount of time: the only difference is that
no update operations are taken, but each node's color is still compared to each
of its neighbor's colors to check that the coloring is valid. Thus each
two-stage iteration is also `O(Vd/N)`, or more precisely `O(Vd/(2N))`.

If we assume that few or no conflicts were found, then we're done. It is fairly
reasonable to assume very few conflicts, since we only have a small number of
cores (e.g., 8 on the test system) but we are looping through tens of thousands
of nodes in our tests. Thus we would expect very few nodes to have a conflict,
and thus the repeated iterations of this algorithm would be performed very
quickly. With these simple assumptions, the entire algorithm is `O(Vd/2N)`. Of
course, these assumptions are false, so it is not ideally parallelizable, as we
will see with the empirical evidence.

If we also make the assumption that the threads have no memory overhead, then
multithreading also offers no extra space complexity (i.e., it is still `O(1)`).
However, it is realistic to say that it is `O(T)`, where `T` is the number of
spawned threads (or goroutines in our case).

##### Empirical results (Revision 1):
Here are some results when running (from the base directory of this repo) on
an i7-2600 (4C8T CPU):
```bash
$ GOPATH=$(pwd) go test -bench=Color* -benchtime=5s -timeout 20m ./src/proj1/
goos: linux
goarch: amd64
pkg: proj1
BenchmarkColorSequentialV100Bf10-8                300000             22569 ns/op
BenchmarkColorSequentialV1000Bf100-8               10000           1099579 ns/op
BenchmarkColorSequentialV10000Bf1000-8               100         101073543 ns/op
BenchmarkColorSequentialV50000Bf5000-8                 2        2616737085 ns/op
BenchmarkColorParallelV100Bf10-8                  100000             93178 ns/op
BenchmarkColorParallelV1000Bf100-8                 10000           1015243 ns/op
BenchmarkColorParallelV10000Bf1000-8                 100          82688240 ns/op
BenchmarkColorParallelV50000Bf5000-8                   3        1810347664 ns/op
```
This includes sequential and parallel runtimes for graphs with uniform average
degree (ignore the number in the middle; this is related to Go's benchmarking
tool). The graph sizes are: G1=(vertices=100, average degree=10); G2=(1000,
100); G3=(10000, 1000); G4=(50000, 5000). The expected runtime, per the analysis
above, is roughly `O(Vd)`, so we would expect G2 to take 2 orders of magnitude
longer than G1, G3 to take 2 orders of magnitude longer than G2, and G4 to take
25 times as long as G3. By our very rough estimate of parallel runtime, we would
expect the time to be reduced by up to a factor of `N/2` (in this case,
`8/2=4`).

We see that this pattern is roughly true for G2, G3, and G4 using the sequential
algorithm (G1 might be too small for it to fit the asymptotic bound).

We see a roughly similar pattern with the parallel runtimes, but we did not
achieve anywhere near an eight-times speedup. There is likely a large overhead
with running the goroutines, especially on smaller graphs (we see that the
sequential algorithm outperforms the parallel one on G1 and almost on G2).
This may be due to a number of things, most likely the creation of so many
goroutines (we create one for each node). It is at least promising that we
achieved any speedup, and it looks like the speedup increases as the graph
gets larger; unfortunately, at this point the test system almost ran out of RAM
and we were unable to experiment further at this point.

Our maximum speedup is 31% on the 50K vertex graph. (We improve this in the next
section.)

In the next revision of the algorithm (see the next section),
we will experiment with lowering the number of goroutines to reduce the memory
overhead and overhead of pseudo-context switches (i.e., the overhead of the
goroutine scheduling), and whatever other optimizations we can find.

##### Changes in Revision 2 and New Results
Structurally, the project was organized a little better with the
introduction of the `proj1/algorithm` package. This holds our sequential
and parallel algorithms from the first revision (`sequential.go` and
`gm_parallel.go`), as well as our newer algorithm from the second revision
(`gm_parallel2.go`).

The second revision differs from the first revision in that:
- The number of goroutines spawned was decreased from `V` to the number
  of specified threads (`nThreads`), which was set to some small multiple
  of the number of logical cores (`runtime.NumCPU() * 2` by default). As
  described earlier, this was an attempt to reduce the computational and memory
  overhead of spawning
  new processes. (This did not have a substantial effect on runtime, which
  goes to show that you can really spam goroutines (at least on the order
  of 10^5) without worrying about much of a performance loss.) We also
  experimented with larger numbers of goroutines, but found that a small number
  was sufficient.
- The buffer used to find the next available valid color for any particular
  node, `neighborColors`, was not reallocated on each loop iteration. Rather,
  we created an empty zero buffer and used `copy` (analogous to C's `memcpy`)
  to fill the array on each iteration. We also implemented this in the
  sequential version. This alone cut all of the algorithm runtimes down roughly
  by a factor of 3 (!!!), both for the sequential and the (new) parallel
  algorithm. (We could not use this on the old parallel algorithm, as each
  goroutine is mapped to a single vertex and thus the `neighborColors` vector
  is not shared over multiple vertex instances.)
- We avoided reallocation of the `U` and `R` vectors, which in practice
  could be very large, by alternating between them. (`U` is the set of vertices
  to be colored in the current iteration, and `R` is the set of vertices to be
  colored in the next round; `U` necessarily subsumes `R`, and we let `U <- R`
  at the end of the iteration, so necessarily `|U1|>=|R1|=|U2|>=|R2|=...>|RN|`.)
  However, this prevents GC on `U` and `R` until the algorithm is finished,
  but this isn't a problem since the memory usage of our algorithm is constant.
- In our first edition, channels (i.e., `chan *graph.Node`) were used to
  build `R` on the mainline thread. While this works fine and fits well within
  the framework of [CSP][csp], we slightly improved performance by switching
  to a mutex lock and directly appending to the array rather than sending
  it to the mainline, since we are basing our algorithm on a shared-memory
  system anyway.
- The `maxColor` variable was used to simplify our algorithm since the
  (approximate) coloring number is not known beforehand. In the first revision,
  we set it to `3*V/2` by default by mistake. We know that the maximum (exact)
  chromatic number is the largest degree of a node + 1; thus, it should be
  `3*d/2`. (The factor of 1.5 is used to ensure safety.) We believe this also
  helped speed up all of the algorithm versions from revision, especially the
  parallel algorithm from revision 1, since:
    - The second bullet point doesn't apply due to the stated reasons
    - This greatly reduces the amount of memory allocation per vertex/goroutine)
    - We didn't actually change any of the implementation of the first
      revision, but `maxColor` is given as a parameter (set in `proj1_test.go`)
      and the benchmark for this unchanged algorithm is much faster than the
      previous revision, especially for larger graphs.

On the same test system as before (i7-2600 (4C/8T) and 8GB RAM), we achieved
the following results. We achieve up to a 48% speedup (2 times speedup) on the
50K vertex graph, which is an improvement over the 31% speedup from the first
revision. Overall, these algorithms run much faster than the previous revision
by virtue of less memory allocation: the revision 1 sequential algorithm takes
4.3 times longer to run than revision 2 sequential, and revision 1 parallel
algorithm takes 5.6 times longer to run than revision 2 parallel (GM2).

```bash
$ GOPATH=$(pwd) go test -bench=Color.* -timeout=20m ./src/proj1
goos: linux
goarch: amd64
pkg: proj1
BenchmarkColorSequentialV100Bf10-8                300000              5815 ns/op
BenchmarkColorSequentialV1000Bf100-8               10000            134939 ns/op
BenchmarkColorSequentialV10000Bf1000-8                50          22013071 ns/op
BenchmarkColorSequentialV50000Bf5000-8                 2         612469096 ns/op
BenchmarkColorParallelGMV100Bf10-8                 20000             80624 ns/op
BenchmarkColorParallelGMV1000Bf100-8                3000            570634 ns/op
BenchmarkColorParallelGMV10000Bf1000-8               100          17103293 ns/op
BenchmarkColorParallelGMV50000Bf5000-8                 3         435580674 ns/op
BenchmarkColorParallelGM2V100Bf10-8                50000             34660 ns/op
BenchmarkColorParallelGM2V1000Bf100-8              10000            139072 ns/op
BenchmarkColorParallelGM2V10000Bf1000-8              100          11396556 ns/op
BenchmarkColorParallelGM2V50000Bf5000-8                5         321242263 ns/op
PASS
ok      proj1   627.908s
```

This time, we were also able to test the program on a test system with a
Ryzen 5 4600H CPU (6C/12T) and 16GB of RAM (and Windows 10), which allowed us
to attempt a benchmark on 100K vertices with an average degree of 5K. This
gives the following results:

```bash
> go test -bench=.*100000.* .\src\proj1
goos: windows
goarch: amd64
pkg: proj1
BenchmarkColorSequentialV100000Bf5000-12               1        4299088700 ns/op
BenchmarkColorParallelGMV100000Bf5000-12               1        1107739000 ns/op
BenchmarkColorParallelGM2V100000Bf5000-12              2         913738400 ns/op
PASS
ok      proj1   81.825s
```

We achieve very close to a 5 times speedup using the new algorithm. (Recall
that the maximum theoretical speedup is `12/2=6` times.) This gives us more
evidence that our algorithm reaches closer to the asymptotic bound as the graph
size increases. It may also just be indicative that more modern CPUs are better
at handling multithreaded applications (R5 is a 2020 CPU, the i7 is a 2011 CPU).
It is possible that we are near our limit at this size, and
that imperfect mechanisms like hyperthreading, context switches, and the
occasional conflict will not let us get much closer to a speedup by a factor
of 6.

##### Next Steps: Scaling Up to Multi-Node
(For project 2)

For project 2, we are aiming to implement [(Gebremedhin 2005)][5], which is
a distributed-memory algorithm based on the ideas in (Gebremedhin 2000)
(the shared-memory variant we used) and Jones and Plassman, among others.

The algorithm from (Gebremedhin 2000) that we used has the downfall that it i
relies on a
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
with 50000 nodes and an average degree of 250 takes over 4GB of RAM;
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
[gam]: http://www.ii.uib.no/~fredrikm/fredrik/papers/Concurrency2000.pdf
[csp]: https://en.wikipedia.org/wiki/Communicating_sequential_processes
