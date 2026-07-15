package similarity

import (
	"math"
	"testing"

	"github.com/tzervas/cowpoke/internal/features"
)

func TestCosineSimilarity(t *testing.T) {
	e := NewEngine(0.85)

	tests := []struct {
		name     string
		v1       []float64
		v2       []float64
		expected float64
	}{
		{
			name:     "identical vectors",
			v1:       []float64{1, 2, 3},
			v2:       []float64{1, 2, 3},
			expected: 1.0,
		},
		{
			name:     "orthogonal vectors",
			v1:       []float64{1, 0, 0},
			v2:       []float64{0, 1, 0},
			expected: 0.0,
		},
		{
			name:     "opposite vectors",
			v1:       []float64{1, 2, 3},
			v2:       []float64{-1, -2, -3},
			expected: -1.0,
		},
		{
			name:     "similar vectors",
			v1:       []float64{1, 2, 3},
			v2:       []float64{1.1, 2.1, 2.9},
			expected: 0.99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := e.CosineSimilarity(tt.v1, tt.v2)
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("CosineSimilarity() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCosineSimilarityEdgeCases(t *testing.T) {
	e := NewEngine(0.85)

	tests := []struct {
		name     string
		v1       []float64
		v2       []float64
		expected float64
	}{
		{
			name:     "different lengths",
			v1:       []float64{1, 2, 3},
			v2:       []float64{1, 2},
			expected: 0.0,
		},
		{
			name:     "empty vectors",
			v1:       []float64{},
			v2:       []float64{},
			expected: 0.0,
		},
		{
			name:     "zero vector",
			v1:       []float64{0, 0, 0},
			v2:       []float64{1, 2, 3},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := e.CosineSimilarity(tt.v1, tt.v2)
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("CosineSimilarity() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMagnitude(t *testing.T) {
	e := NewEngine(0.85)

	tests := []struct {
		name     string
		vector   []float64
		expected float64
	}{
		{"unit vector", []float64{1, 0, 0}, 1.0},
		{"3-4-5 triangle", []float64{3, 4}, 5.0},
		{"zero vector", []float64{0, 0, 0}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := e.magnitude(tt.vector)
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("magnitude() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	e := NewEngine(0.85)

	vector := []float64{3, 4}
	normalized := e.Normalize(vector)

	// Check magnitude is 1
	mag := e.magnitude(normalized)
	if math.Abs(mag-1.0) > 0.01 {
		t.Errorf("Normalized vector magnitude = %v, want 1.0", mag)
	}

	// Check direction preserved
	expected := []float64{0.6, 0.8}
	for i := range normalized {
		if math.Abs(normalized[i]-expected[i]) > 0.01 {
			t.Errorf("Normalize()[%d] = %v, want %v", i, normalized[i], expected[i])
		}
	}
}

func TestEuclideanDistance(t *testing.T) {
	e := NewEngine(0.85)

	tests := []struct {
		name     string
		v1       []float64
		v2       []float64
		expected float64
		wantErr  bool
	}{
		{"identical", []float64{1, 2, 3}, []float64{1, 2, 3}, 0.0, false},
		{"3-4-5 triangle", []float64{0, 0}, []float64{3, 4}, 5.0, false},
		{"different lengths", []float64{1, 2}, []float64{1, 2, 3}, 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.EuclideanDistance(tt.v1, tt.v2)
			if (err != nil) != tt.wantErr {
				t.Errorf("EuclideanDistance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("EuclideanDistance() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestComputeSimilarity(t *testing.T) {
	e := NewEngine(0.85)

	fv1 := &features.FeatureVector{
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

	fv2 := &features.FeatureVector{
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

	similarity := e.ComputeSimilarity(fv1, fv2)
	if math.Abs(similarity-1.0) > 0.01 {
		t.Errorf("ComputeSimilarity() = %v, want 1.0", similarity)
	}
}

func TestIsSimilar(t *testing.T) {
	e := NewEngine(0.85)

	fv1 := &features.FeatureVector{
		Mean: 50.0, CV: 0.1, P50: 49.0, P95: 55.0, P99: 57.0,
		Skewness: 0.5, Kurtosis: 0.3, Slope: 0.2, Autocorr: 0.8,
	}

	fv2 := &features.FeatureVector{
		Mean: 50.0, CV: 0.1, P50: 49.0, P95: 55.0, P99: 57.0,
		Skewness: 0.5, Kurtosis: 0.3, Slope: 0.2, Autocorr: 0.8,
	}

	if !e.IsSimilar(fv1, fv2) {
		t.Error("IsSimilar() should return true for identical vectors")
	}
}

func TestThreshold(t *testing.T) {
	e := NewEngine(0.85)

	if e.GetThreshold() != 0.85 {
		t.Errorf("GetThreshold() = %v, want 0.85", e.GetThreshold())
	}

	e.SetThreshold(0.90)
	if e.GetThreshold() != 0.90 {
		t.Errorf("GetThreshold() after SetThreshold() = %v, want 0.90", e.GetThreshold())
	}
}
