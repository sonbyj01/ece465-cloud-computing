# ece465-cloud-computing
ECE465 Cloud Computing Homework Assignments

Jonathan Lam & Henry Son

Prof. Marano

---

We chose to focus on the graph coloring algorithm for projects 1 and 2.

- [Project 1](./src/proj1/README.md): Graph coloring in a single-node,
  multithreaded environment
- Project 2: Graph coloring in a multi-node, multithreaded environment 

---

### Build Instructions
Make sure Golang is installed on your system. The test system (for project 1)
used `go1.11.6 linux/amd64`.

From the current directory, run:
```bash
GOPATH=$(pwd) go run src/projX/projX.go
```
where `X` is the project number. For more information about each project,
see the README in `src/projX/` folder.

---

### Tests/Benchmarks
Benchmarks are provided in `src/projX/projX_test.go`. To run tests:
```bash
GOPATH=$(pwd) go test ./src/projX
```
To run benchmarks:
```bash
GOPATH=$(pwd) go test -bench=. ./src/projX
```
To run a specific benchmark, put the name after the `-bench`. For example, to
run the benchmarks for generating new random graphs in `src/proj1/proj1_test.go`
(the `BenchmarkNewRandomGraph` and `BenchmarkNewRandomGraphParallel` functions):
```bash
GOPATH=$(pwd) go test -bench=NewRandomGraph* ./src/projX
```
See `src/projX/projX_test.go` for available tests and benchmarks. See the
[Golang documentation](https://golang.org/doc/) for more details on the
build/run environment.