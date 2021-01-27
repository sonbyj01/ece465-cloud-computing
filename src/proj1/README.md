# ECE465 Project 1: Graph Coloring

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

While finding the (true) chromatic number is an
NP-complete problem (and thus the only exact algorithms are
exponential-time algorithms), there exist efficient randomized, distributed
algorithms that can approximate the chromatic number (and perform fairly
well in real-world tasks). The algorithm we choose is called the
Gebremedhin-Manne algorithm, which is described in [this presentation][1];
this takes a greedy approach by speculatively choosing the lowest
whole number (color) not taken by any of its neighbors, for multiple
nodes in parallel. Since this step is done without any data synchronization,
it is fast but may result in some inconsistencies (conflicts); thus
after performing this step on every node, we perform the same procedure
again only on the nodes which had conflicts, repeating until no conflicts
are found.

### Resources

- [A greedy parallel algorithm for Δ+1 graph coloring
  (simple but bad estimate for chromatic number)][2]
- [A smaller exponential time bound for exact graph coloring (using
  small maximal independent sets)][4]
- [The Gebremedhin-Manne algorithm: A randomized, distributed,
  speculative-coloring approach to graph coloring][5]
- [Distributed algorithm for finding maximal independent sets (MIS)][3]
    - Luby's algorithm: another distributed algorithm for MIS
    - Jones and Plassman's method: a graph coloring method using MIS
- [Graph coloring on the GPU (with a good literature review of distributed
  graph coloring algorithms)][6]

[1]: https://www.osti.gov/servlets/purl/1246285
[2]: https://stanford.edu/~rezab/classes/cme323/S16/projects_reports/bae.pdf
[3]: https://tik-old.ee.ethz.ch/file//2be1291694b1730bba83f7fa18d9e0f2/podc08SW.pdf
[4]: https://arxiv.org/pdf/cs/0011009.pdf
[5]: https://cscapes.cs.purdue.edu/coloringpage/abstracts/euro05.pdf
[6]: https://people.eecs.berkeley.edu/~aydin/coloring.pdf 