# ece465-cloud-computing
ECE465 Cloud Computing Homework Assignments

Jonathan Lam & Henry Son

Prof. Marano

---

### Build Instructions
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

TODO: elaborate here