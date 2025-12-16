package similarity

import (
	"fmt"
	"math"

	"github.com/tzervas/cowpoke/internal/features"
)

// Engine computes similarity between feature vectors
type Engine struct {
	threshold float64
}

// NewEngine creates a new similarity engine
func NewEngine(threshold float64) *Engine {
	return &Engine{
		threshold: threshold,
	}
}

// ComputeSimilarity calculates cosine similarity between two feature vectors
func (e *Engine) ComputeSimilarity(v1, v2 *features.FeatureVector) float64 {
	vec1 := v1.ToVector()
	vec2 := v2.ToVector()

	return e.CosineSimilarity(vec1, vec2)
}

// CosineSimilarity computes the cosine similarity between two vectors
// Returns a value between -1 and 1, where:
//
//	 1 = identical direction
//	 0 = orthogonal (no similarity)
//	-1 = opposite direction
func (e *Engine) CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0.0
	}

	if len(a) == 0 {
		return 0.0
	}

	dotProduct := e.dotProduct(a, b)
	magnitudeA := e.magnitude(a)
	magnitudeB := e.magnitude(b)

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0.0
	}

	return dotProduct / (magnitudeA * magnitudeB)
}

// dotProduct computes the dot product of two vectors
func (e *Engine) dotProduct(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

// magnitude computes the magnitude (L2 norm) of a vector
func (e *Engine) magnitude(v []float64) float64 {
	sum := 0.0
	for _, val := range v {
		sum += val * val
	}
	return math.Sqrt(sum)
}

// EuclideanDistance computes the Euclidean distance between two vectors
func (e *Engine) EuclideanDistance(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		return 0.0, fmt.Errorf("vectors must have same length")
	}

	sum := 0.0
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return math.Sqrt(sum), nil
}

// Normalize normalizes a vector to unit length
func (e *Engine) Normalize(v []float64) []float64 {
	mag := e.magnitude(v)
	if mag == 0 {
		return v
	}

	normalized := make([]float64, len(v))
	for i, val := range v {
		normalized[i] = val / mag
	}
	return normalized
}

// IsSimilar checks if similarity exceeds the threshold
func (e *Engine) IsSimilar(v1, v2 *features.FeatureVector) bool {
	similarity := e.ComputeSimilarity(v1, v2)
	return similarity >= e.threshold
}

// SetThreshold updates the similarity threshold
func (e *Engine) SetThreshold(threshold float64) {
	e.threshold = threshold
}

// GetThreshold returns the current similarity threshold
func (e *Engine) GetThreshold() float64 {
	return e.threshold
}
