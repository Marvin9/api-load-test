package pkg

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/montanaflynn/stats"
)

type AnalyzedData struct {
	Mean              float64
	Median            float64
	Mode              []float64
	Range             float64
	StandardDeviation float64
	Min               float64
	Max               float64
}

func (a AnalyzedData) Display() {
	modes := ""
	for _, mode := range a.Mode {
		modes += fmt.Sprintf("%v, ", mode)
	}
	modes = strings.TrimSuffix(modes, ", ")
	fmt.Printf(`
Mean: %vms
Median: %vms
Mode: %vms
Range: %vms
Standard Deviation: %vms
Min: %vms
Max: %vms
`, a.Mean, a.Median, modes, a.Range, a.StandardDeviation, a.Min, a.Max)
}

func checkCalcError(operation string, err error) {
	if err != nil {
		log.Fatalf("\n[ERROR] calculation %v.\n%v", operation, err)
	}
}

func Analysis(report Report) AnalyzedData {
	report = report[1:]
	performance := make([]float64, len(report))
	for i, data := range report {
		performance[i] = float64(data.Performance)
	}

	sort.Float64s(performance)

	mean, err := stats.Mean(performance)
	checkCalcError("mean", err)

	min := performance[0]
	max := performance[len(performance)-1]

	median, err := stats.Median(performance)
	checkCalcError("median", err)

	mode, err := stats.Mode(performance)
	checkCalcError("mode", err)

	stdDev, err := stats.StandardDeviation(performance)
	checkCalcError("standard deviation", err)

	return AnalyzedData{
		Mean:              mean,
		Median:            median,
		Mode:              mode,
		Range:             max - min,
		StandardDeviation: stdDev,
		Min:               min,
		Max:               max,
	}
}
