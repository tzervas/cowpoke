package features

import (
	"fmt"
	"math"
	"sort"
)

// Extractor extracts statistical features from metric time series
type Extractor struct {
	windowSize int
}

// FeatureVector represents extracted features from metrics
type FeatureVector struct {
	Mean     float64
	CV       float64 // Coefficient of Variation
	P50      float64 // 50th percentile (median)
	P95      float64 // 95th percentile
	P99      float64 // 99th percentile
	Skewness float64
	Kurtosis float64
	Slope    float64 // Linear trend slope
	Autocorr float64 // Lag-1 autocorrelation
}

// NewExtractor creates a new feature extractor
func NewExtractor(windowSize int) *Extractor {
	return &Extractor{
		windowSize: windowSize,
	}
}

// Extract computes all statistical features from the given metrics
func (e *Extractor) Extract(metrics []float64) (*FeatureVector, error) {
	if len(metrics) == 0 {
		return nil, fmt.Errorf("no metrics provided")
	}

	fv := &FeatureVector{}

	// Compute basic statistics
	fv.Mean = e.computeMean(metrics)
	stdDev := e.computeStdDev(metrics, fv.Mean)

	// Coefficient of Variation (CV = stdDev / mean)
	if fv.Mean != 0 {
		fv.CV = stdDev / fv.Mean
	}

	// Percentiles
	fv.P50 = e.computePercentile(metrics, 0.50)
	fv.P95 = e.computePercentile(metrics, 0.95)
	fv.P99 = e.computePercentile(metrics, 0.99)

	// Higher-order moments
	fv.Skewness = e.computeSkewness(metrics, fv.Mean, stdDev)
	fv.Kurtosis = e.computeKurtosis(metrics, fv.Mean, stdDev)

	// Temporal features
	fv.Slope = e.computeSlope(metrics)
	fv.Autocorr = e.computeAutocorrelation(metrics, fv.Mean)

	return fv, nil
}

// computeMean calculates the arithmetic mean
func (e *Extractor) computeMean(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

// computeStdDev calculates the standard deviation
func (e *Extractor) computeStdDev(data []float64, mean float64) float64 {
	if len(data) <= 1 {
		return 0.0
	}

	variance := 0.0
	for _, v := range data {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(data) - 1)
	return math.Sqrt(variance)
}

// computePercentile calculates the specified percentile
func (e *Extractor) computePercentile(data []float64, p float64) float64 {
	if len(data) == 0 {
		return 0.0
	}

	// Create a copy to avoid modifying original
	sorted := make([]float64, len(data))
	copy(sorted, data)
	sort.Float64s(sorted)

	index := p * float64(len(sorted)-1)
	lower := int(math.Floor(index))
	upper := int(math.Ceil(index))

	if lower == upper {
		return sorted[lower]
	}

	// Linear interpolation
	weight := index - float64(lower)
	return sorted[lower]*(1-weight) + sorted[upper]*weight
}

// computeSkewness calculates the skewness (Fisher-Pearson coefficient)
func (e *Extractor) computeSkewness(data []float64, mean, stdDev float64) float64 {
	if stdDev == 0 || len(data) < 3 {
		return 0.0
	}

	n := float64(len(data))
	sum := 0.0
	for _, v := range data {
		z := (v - mean) / stdDev
		sum += z * z * z
	}

	// Sample skewness with bias correction
	skew := sum / n
	correction := math.Sqrt(n*(n-1)) / (n - 2)
	return skew * correction
}

// computeKurtosis calculates the excess kurtosis
func (e *Extractor) computeKurtosis(data []float64, mean, stdDev float64) float64 {
	if stdDev == 0 || len(data) < 4 {
		return 0.0
	}

	n := float64(len(data))
	sum := 0.0
	for _, v := range data {
		z := (v - mean) / stdDev
		sum += z * z * z * z
	}

	// Excess kurtosis (subtract 3 for normal distribution baseline)
	kurt := (sum / n) - 3.0
	return kurt
}

// computeSlope calculates the linear trend slope using least squares
func (e *Extractor) computeSlope(data []float64) float64 {
	n := float64(len(data))
	if n < 2 {
		return 0.0
	}

	// Use index as x-axis (time)
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0

	for i, y := range data {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	// Slope = (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	denominator := n*sumX2 - sumX*sumX
	if denominator == 0 {
		return 0.0
	}

	slope := (n*sumXY - sumX*sumY) / denominator
	return slope
}

// computeAutocorrelation calculates lag-1 autocorrelation
func (e *Extractor) computeAutocorrelation(data []float64, mean float64) float64 {
	n := len(data)
	if n < 2 {
		return 0.0
	}

	// Compute variance
	variance := 0.0
	for _, v := range data {
		diff := v - mean
		variance += diff * diff
	}

	if variance == 0 {
		return 0.0
	}

	// Compute lag-1 autocovariance
	autocovariance := 0.0
	for i := 0; i < n-1; i++ {
		autocovariance += (data[i] - mean) * (data[i+1] - mean)
	}

	// Autocorrelation = autocovariance / variance
	return autocovariance / variance
}

// ToVector converts FeatureVector to a slice for similarity computation
func (fv *FeatureVector) ToVector() []float64 {
	return []float64{
		fv.Mean,
		fv.CV,
		fv.P50,
		fv.P95,
		fv.P99,
		fv.Skewness,
		fv.Kurtosis,
		fv.Slope,
		fv.Autocorr,
	}
}
