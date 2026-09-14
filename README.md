# Collision Evaluation Study

Empirical comparative analysis of spatial data structures for broad-phase collision detection under high entity churn (1,000–20,000 moving entities).

## Core Research Question
What are the exact performance trade-offs between dynamic tree-rebuilding overhead and collision traversal efficiency under continuous entity motion?

Comparing:
- **Bounding Volume Hierarchies (BVH)** with incremental refitting
- **K-Dimensional Trees (KD-Trees)** with spatial rebuilding
- **Brute-Force $O(N^2)$ Baseline** as correctness oracle and baseline floor

## Project Structure
```text
.
├── Makefile
├── cmd/
│   └── bench/          # CLI benchmark runner
├── pkg/
│   ├── baseline/       # Brute-force O(N²) reference implementation
│   ├── geom/           # 2D primitives (Vec2, AABB2)
│   └── sim/            # Entity, World, deterministic tick loop, metrics/CSV exporter
```

## Getting Started

### Run Tests
```bash
make test
# or: go test ./... -v -count=1
```

### Run Benchmark Simulation
```bash
make run
# or custom parameters:
go run ./cmd/bench/ -n 1000 -frames 100 -seed 42 -out results.csv
```
