package sim

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"sort"
)

// FrameRecord holds per-frame telemetry.
type FrameRecord struct {
	EntityCount    int
	TreeUpdateMs   float64
	QueryMs        float64
	TotalFrameMs   float64
	HeapAllocBytes uint64
	GcPauseMs      float64
	CandidatePairs int
}

// Recorder accumulates per-frame records.
type Recorder struct {
	Frames []FrameRecord
}

// Record appends a frame record.
func (r *Recorder) Record(f FrameRecord) {
	r.Frames = append(r.Frames, f)
}

// Percentiles computes p50, p95, p99 of TotalFrameMs across all recorded frames.
// Returns 0,0,0 if no frames are recorded.
func (r *Recorder) Percentiles() (p50, p95, p99 float64) {
	n := len(r.Frames)
	if n == 0 {
		return 0, 0, 0
	}
	vals := make([]float64, n)
	for i, f := range r.Frames {
		vals[i] = f.TotalFrameMs
	}
	sort.Float64s(vals)
	p50 = vals[pctIndex(n, 50)]
	p95 = vals[pctIndex(n, 95)]
	p99 = vals[pctIndex(n, 99)]
	return
}

func pctIndex(n, pct int) int {
	idx := int(math.Ceil(float64(pct)/100*float64(n))) - 1
	if idx < 0 {
		return 0
	}
	if idx >= n {
		return n - 1
	}
	return idx
}

// WriteCSV writes all frame records as CSV to w.
func (r *Recorder) WriteCSV(w io.Writer) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	header := []string{
		"entity_count", "tree_update_ms", "query_ms", "total_frame_ms",
		"heap_alloc_bytes", "gc_pause_ms", "candidate_pairs",
	}
	if err := cw.Write(header); err != nil {
		return err
	}

	for _, f := range r.Frames {
		row := []string{
			fmt.Sprintf("%d", f.EntityCount),
			fmt.Sprintf("%.4f", f.TreeUpdateMs),
			fmt.Sprintf("%.4f", f.QueryMs),
			fmt.Sprintf("%.4f", f.TotalFrameMs),
			fmt.Sprintf("%d", f.HeapAllocBytes),
			fmt.Sprintf("%.4f", f.GcPauseMs),
			fmt.Sprintf("%d", f.CandidatePairs),
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	return cw.Error()
}
