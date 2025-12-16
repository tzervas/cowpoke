package features

import (
	"math"
	"testing"
)

func TestExtractorMean(t *testing.T) {
	e := NewExtractor(60)

	tests := []struct {
		name     string
		metrics  []float64
		expected float64
	}{
		{"simple", []float64{1, 2, 3, 4, 5}, 3.0},
		{"single", []float64{10}, 10.0},
		{"zeros", []float64{0, 0, 0}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := e.computeMean(tt.metrics)
			if math.Abs(result-tt.expected) > 0.001 {
				t.Errorf("computeMean() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractorPercentile(t *testing.T) {
	e := NewExtractor(60)
	metrics := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	tests := []struct {
		name       string
		percentile float64
		expected   float64
	}{
		{"median", 0.5, 5.5},
		{"p95", 0.95, 9.55},
		{"p99", 0.99, 9.91},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := e.computePercentile(metrics, tt.percentile)
			if math.Abs(result-tt.expected) > 0.1 {
				t.Errorf("computePercentile(%v) = %v, want %v", tt.percentile, result, tt.expected)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	e := NewExtractor(60)
	metrics := []float64{50, 52, 48, 55, 51, 49, 53, 50}

	fv, err := e.Extract(metrics)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if fv == nil {
		t.Fatal("Extract() returned nil feature vector")
	}

	// Basic sanity checks
	if fv.Mean <= 0 {
		t.Errorf("Mean should be positive, got %v", fv.Mean)
	}

	if fv.CV < 0 {
		t.Errorf("CV should be non-negative, got %v", fv.CV)
	}

	if fv.P50 <= 0 {
		t.Errorf("P50 should be positive, got %v", fv.P50)
	}
}

func TestExtractEmptyMetrics(t *testing.T) {
	e := NewExtractor(60)
	metrics := []float64{}

	_, err := e.Extract(metrics)
	if err == nil {
		t.Error("Extract() should return error for empty metrics")
	}
}

func TestFeatureVectorToVector(t *testing.T) {
	fv := &FeatureVector{
		Mean:     50.0,
		CV:       0.1,
		P50:      49.0,
		P95:      55.0,
		P99:      57.0,
		Skewness: 0.5,
		Kurtosis: 0.3,
		Slope:    0.2,
		Autocorr: 0.8,
	}

	vec := fv.ToVector()
	if len(vec) != 9 {
		t.Errorf("ToVector() should return 9 elements, got %d", len(vec))
	}

	if vec[0] != fv.Mean {
		t.Errorf("First element should be Mean, got %v", vec[0])
	}
}
