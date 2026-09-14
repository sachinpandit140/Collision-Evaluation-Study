package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/sachinpandit140/Collision-Evaluation-Study/pkg/baseline"
	"github.com/sachinpandit140/Collision-Evaluation-Study/pkg/sim"
)

func main() {
	n := flag.Int("n", 1000, "number of entities")
	frames := flag.Int("frames", 100, "number of simulation frames")
	seed := flag.Int64("seed", 42, "random seed for determinism")
	out := flag.String("out", "results.csv", "CSV output path")
	flag.Parse()

	fmt.Printf("Running brute-force baseline: %d entities, %d frames, seed=%d\n", *n, *frames, *seed)

	w := sim.NewWorld(1000, 1000, *n, *seed)
	rec := &sim.Recorder{}
	dt := 1.0

	for f := 0; f < *frames; f++ {
		var mem runtime.MemStats

		frameStart := time.Now()

		// Tree update phase (brute-force has no tree, so this is just the tick)
		tickStart := time.Now()
		w.Tick(dt)
		tickMs := float64(time.Since(tickStart).Microseconds()) / 1000.0

		// Query phase (brute-force collision check)
		runtime.ReadMemStats(&mem)
		allocBefore := mem.TotalAlloc

		queryStart := time.Now()
		pairs := baseline.FindCollisions(w.Entities)
		queryMs := float64(time.Since(queryStart).Microseconds()) / 1000.0

		runtime.ReadMemStats(&mem)
		allocDelta := mem.TotalAlloc - allocBefore

		totalMs := float64(time.Since(frameStart).Microseconds()) / 1000.0

		rec.Record(sim.FrameRecord{
			EntityCount:    len(w.Entities),
			TreeUpdateMs:   tickMs,
			QueryMs:        queryMs,
			TotalFrameMs:   totalMs,
			HeapAllocBytes: allocDelta,
			GcPauseMs:      0, // TODO: extract from runtime.MemStats.PauseNs
			CandidatePairs: len(pairs),
		})
	}

	// Write CSV
	f, err := os.Create(*out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating %s: %v\n", *out, err)
		os.Exit(1)
	}
	defer f.Close()

	if err := rec.WriteCSV(f); err != nil {
		fmt.Fprintf(os.Stderr, "error writing CSV: %v\n", err)
		os.Exit(1)
	}

	// Print summary
	p50, p95, p99 := rec.Percentiles()
	fmt.Printf("\nResults written to %s (%d frames)\n", *out, len(rec.Frames))
	fmt.Printf("  p50: %.3f ms\n", p50)
	fmt.Printf("  p95: %.3f ms\n", p95)
	fmt.Printf("  p99: %.3f ms\n", p99)
}
