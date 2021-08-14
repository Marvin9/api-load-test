package pkg

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/montanaflynn/stats"
)

type AnalyzedData struct {
	Mean                  float64
	Median                float64
	Mode                  []float64
	Range                 float64
	StandardDeviation     float64
	Min                   float64
	Max                   float64
	NintiethPercentaile   float64
	NintyFifthPercentaile float64
}

func (a AnalyzedData) Display() {
	modes := ""
	for _, mode := range a.Mode {
		modes += fmt.Sprintf("%.3f, ", mode)
	}
	modes = strings.TrimSuffix(modes, ", ")
	fmt.Printf(`
	Summary
	----------------------------------
	| Mean:               %.3fms
	| Median:             %.3fms
	| Mode:               %vms
	| Range:              %.3fms
	| Standard Deviation: %.3fms
	| Min:                %.3fms
	| Max:                %.3fms
	| 90th Percentaile:   %.3fms
	| 95th Percentaile:   %.3fms
	---------------------------------

`, a.Mean, a.Median, modes, a.Range, a.StandardDeviation, a.Min, a.Max, a.NintiethPercentaile, a.NintyFifthPercentaile)
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

	nintiethPercentaile, err := stats.Percentile(performance, float64(90))
	checkCalcError("90th percentaile", err)

	nintyFifthPercentaile, err := stats.Percentile(performance, float64(95))
	checkCalcError("95th percentaile", err)

	return AnalyzedData{
		Mean:                  mean,
		Median:                median,
		Mode:                  mode,
		Range:                 max - min,
		StandardDeviation:     stdDev,
		Min:                   min,
		Max:                   max,
		NintiethPercentaile:   nintiethPercentaile,
		NintyFifthPercentaile: nintyFifthPercentaile,
	}
}
